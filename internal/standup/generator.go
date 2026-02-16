package standup

import (
	"context"
	"fmt"
	"path/filepath"
	"strings"
	"time"

	"github.com/kushalpatel/bro-cli/internal/git"
	"github.com/kushalpatel/bro-cli/internal/llm"
)

// Generator handles scanning repos and producing standup summaries.
type Generator struct {
	Repos    []string
	Provider llm.Provider
}

// NewGenerator creates a standup generator for the given repo paths.
func NewGenerator(repos []string, provider llm.Provider) *Generator {
	return &Generator{Repos: repos, Provider: provider}
}

// SinceDate returns the appropriate "since" date for standup.
// On Monday, it returns last Friday. Otherwise, yesterday.
func SinceDate() string {
	now := time.Now()
	var since time.Time
	switch now.Weekday() {
	case time.Monday:
		since = now.AddDate(0, 0, -3) // last Friday
	default:
		since = now.AddDate(0, 0, -1) // yesterday
	}
	return since.Format("2006-01-02")
}

// Generate scans repos for recent git activity and produces an LLM summary.
func (g *Generator) Generate() (string, error) {
	since := SinceDate()
	var allLogs strings.Builder

	for _, repo := range g.Repos {
		log, err := git.Log(repo, since)
		if err != nil {
			fmt.Fprintf(&allLogs, "\n## %s\n(error reading repo: %v)\n", filepath.Base(repo), err)
			continue
		}
		if log == "" {
			continue
		}
		fmt.Fprintf(&allLogs, "\n## %s\n%s\n", filepath.Base(repo), log)
	}

	if allLogs.Len() == 0 {
		return "No git activity found since " + since + ".", nil
	}

	systemPrompt := `You are a developer writing a standup summary. Summarize the git activity into a concise standup update with these sections:
- **Yesterday**: What was accomplished
- **Today**: What's planned (infer from context or say "Continue work on..." based on recent activity)
- **Blockers**: Any potential blockers (say "None" if unclear)

Be concise. Use bullet points. Don't include commit hashes.`

	prompt := fmt.Sprintf("Here are my git commits since %s:\n%s\n\nWrite my standup summary.", since, allLogs.String())

	ctx := context.Background()
	summary, err := g.Provider.Generate(ctx, prompt, systemPrompt)
	if err != nil {
		return "", fmt.Errorf("generate standup summary: %w", err)
	}

	return summary, nil
}
