package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/kushalpatel/bro-cli/internal/config"
	"github.com/kushalpatel/bro-cli/internal/llm"
	"github.com/kushalpatel/bro-cli/internal/standup"
	"github.com/kushalpatel/bro-cli/internal/ui"
	"github.com/spf13/cobra"
)

var standupCmd = &cobra.Command{
	Use:   "standup",
	Short: "Generate a standup summary from git logs",
	Long:  "Scans configured repos for recent git activity and generates an LLM-summarized standup report.",
	RunE:  runStandup,
}

func init() {
	rootCmd.AddCommand(standupCmd)
}

func runStandup(cmd *cobra.Command, args []string) error {
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("load config: %w", err)
	}

	if err := ensureSetup(cfg); err != nil {
		return err
	}

	if len(cfg.StandupRepos) == 0 {
		return fmt.Errorf("no repos configured for standup. Add 'standup_repos' to ~/.bro/config.yaml")
	}

	provider, err := llm.NewProvider(cfg)
	if err != nil {
		return err
	}

	gen := standup.NewGenerator(cfg.StandupRepos, provider)

	fmt.Println(ui.DimStyle.Render("Scanning repos and generating standup..."))

	summary, err := gen.Generate()
	if err != nil {
		return fmt.Errorf("generate standup: %w", err)
	}

	fmt.Println()
	fmt.Println(summary)

	// Save to file
	configDir, err := config.Dir()
	if err != nil {
		return err
	}
	if err := config.EnsureDir(); err != nil {
		return err
	}

	filename := time.Now().Format("2006-01-02") + ".md"
	savePath := filepath.Join(configDir, "standups", filename)

	if err := os.WriteFile(savePath, []byte(summary+"\n"), 0644); err != nil {
		return fmt.Errorf("save standup: %w", err)
	}

	fmt.Println()
	fmt.Println(ui.SuccessStyle.Render("Saved to " + savePath))

	return nil
}
