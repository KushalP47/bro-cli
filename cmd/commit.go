package cmd

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/kushalpatel/bro-cli/internal/config"
	"github.com/kushalpatel/bro-cli/internal/git"
	"github.com/kushalpatel/bro-cli/internal/llm"
	"github.com/kushalpatel/bro-cli/internal/ui"
	"github.com/spf13/cobra"
)

var commitCmd = &cobra.Command{
	Use:   "commit",
	Short: "Generate a commit message from staged changes",
	Long:  "Analyzes git diff --staged and uses an LLM to generate a commit message. Supports accept, edit, and regenerate flow.",
	RunE:  runCommit,
}

func init() {
	rootCmd.AddCommand(commitCmd)
}

const maxDiffLen = 8000

func runCommit(cmd *cobra.Command, args []string) error {
	diff, err := git.StagedDiff()
	if err != nil {
		return fmt.Errorf("get staged diff: %w", err)
	}
	if diff == "" {
		return fmt.Errorf("no staged changes. Use 'git add' to stage files first")
	}

	// Truncate large diffs
	if len(diff) > maxDiffLen {
		diff = diff[:maxDiffLen] + "\n... (truncated)"
	}

	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("load config: %w", err)
	}

	// Try to create LLM provider
	var provider llm.Provider
	if cfg.Provider != "" || cfg.Anthropic.APIKey != "" || cfg.OpenAI.APIKey != "" ||
		os.Getenv("BRO_ANTHROPIC_API_KEY") != "" || os.Getenv("BRO_OPENAI_API_KEY") != "" {
		provider, err = llm.NewProvider(cfg)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Warning: LLM unavailable (%v). Falling back to manual mode.\n", err)
		}
	}

	// Match profile
	cwd, _ := os.Getwd()
	profile := config.MatchProfile(cfg.Profiles, cwd)

	if provider == nil {
		return fallbackCommit(profile)
	}

	// Generate message with LLM
	systemPrompt, userPrompt := buildCommitPrompts(diff, profile)

	for {
		ctx := context.Background()
		message, err := provider.Generate(ctx, userPrompt, systemPrompt)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Warning: LLM error (%v). Falling back to manual mode.\n", err)
			return fallbackCommit(profile)
		}

		// If template profile, apply template
		if profile != nil && profile.CommitStyle == "template" && profile.Template != "" {
			message = applyCommitTemplate(message, profile)
		}

		choice, err := ui.RunCommitFlow(message)
		if err != nil {
			return fmt.Errorf("commit flow: %w", err)
		}

		switch choice {
		case ui.CommitAccept:
			return execCommit(message)
		case ui.CommitEdit:
			edited, err := editMessage(message)
			if err != nil {
				return err
			}
			if edited != "" {
				return execCommit(edited)
			}
			return nil
		case ui.CommitRegenerate:
			continue
		case ui.CommitQuit:
			fmt.Println("Commit cancelled.")
			return nil
		}
	}
}

func buildCommitPrompts(diff string, profile *config.Profile) (string, string) {
	if profile != nil && profile.CommitStyle == "conventional" {
		systemPrompt := `You are a commit message generator. Generate a concise conventional commit message.
Format: <type>(<scope>): <description>
Types: feat, fix, chore, docs, refactor, test, ci
Rules:
- Use lowercase
- No period at the end
- Description should be imperative mood ("add" not "added")
- Keep under 72 characters
- Only output the commit message, nothing else`

		userPrompt := fmt.Sprintf("Generate a conventional commit message for this diff:\n\n%s", diff)
		return systemPrompt, userPrompt
	}

	if profile != nil && profile.CommitStyle == "template" {
		systemPrompt := `You are a commit message generator. Generate ONLY a concise commit message body.
Rules:
- Imperative mood ("add" not "added")
- Keep under 60 characters
- No prefix, no type, just the message body
- Only output the message, nothing else`

		userPrompt := fmt.Sprintf("Generate a commit message body for this diff:\n\n%s", diff)
		return systemPrompt, userPrompt
	}

	// Default: conventional commits
	systemPrompt := `You are a commit message generator. Generate a concise conventional commit message.
Format: <type>: <description>
Types: feat, fix, chore, docs, refactor, test, ci
Rules:
- Use lowercase
- No period at the end
- Imperative mood
- Keep under 72 characters
- Only output the commit message, nothing else`

	userPrompt := fmt.Sprintf("Generate a commit message for this diff:\n\n%s", diff)
	return systemPrompt, userPrompt
}

func applyCommitTemplate(message string, profile *config.Profile) string {
	branch, err := git.CurrentBranch()
	if err != nil {
		return message
	}

	if profile.BranchPattern == "" {
		return message
	}

	parts, err := git.ParseBranch(branch, profile.BranchPattern)
	if err != nil {
		return message
	}

	result := profile.Template
	result = strings.ReplaceAll(result, "{message}", message)
	for key, value := range parts {
		result = strings.ReplaceAll(result, "{"+key+"}", value)
	}
	return result
}

func execCommit(message string) error {
	gitCmd := exec.Command("git", "commit", "-m", message)
	gitCmd.Stdout = os.Stdout
	gitCmd.Stderr = os.Stderr
	if err := gitCmd.Run(); err != nil {
		return fmt.Errorf("git commit: %w", err)
	}
	fmt.Println(ui.SuccessStyle.Render("Committed!"))
	return nil
}

func editMessage(message string) (string, error) {
	editor := os.Getenv("EDITOR")
	if editor == "" {
		editor = "vi"
	}

	tmpFile, err := os.CreateTemp("", "bro-commit-*.txt")
	if err != nil {
		return "", fmt.Errorf("create temp file: %w", err)
	}
	defer os.Remove(tmpFile.Name())

	if _, err := tmpFile.WriteString(message); err != nil {
		return "", fmt.Errorf("write temp file: %w", err)
	}
	tmpFile.Close()

	editorCmd := exec.Command(editor, tmpFile.Name())
	editorCmd.Stdin = os.Stdin
	editorCmd.Stdout = os.Stdout
	editorCmd.Stderr = os.Stderr
	if err := editorCmd.Run(); err != nil {
		return "", fmt.Errorf("run editor: %w", err)
	}

	data, err := os.ReadFile(tmpFile.Name())
	if err != nil {
		return "", fmt.Errorf("read edited file: %w", err)
	}

	return strings.TrimSpace(string(data)), nil
}

func fallbackCommit(profile *config.Profile) error {
	var template string
	if profile != nil && profile.Template != "" {
		template = profile.Template
	} else {
		template = "feat: "
	}

	edited, err := editMessage(template)
	if err != nil {
		return err
	}
	if edited == "" || edited == template {
		fmt.Println("Commit cancelled (empty message).")
		return nil
	}
	return execCommit(edited)
}
