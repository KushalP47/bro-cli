package cmd

import (
	"context"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/kushalpatel/bro-cli/internal/config"
	"github.com/kushalpatel/bro-cli/internal/llm"
	"github.com/spf13/cobra"
)

var askCmd = &cobra.Command{
	Use:   "ask [question]",
	Short: "Ask the LLM a question",
	Long:  "General LLM queries with support for stdin pipe and -f file flag. Streams response to terminal.",
	RunE:  runAsk,
}

func init() {
	askCmd.Flags().StringP("file", "f", "", "File to include as context")
	rootCmd.AddCommand(askCmd)
}

func runAsk(cmd *cobra.Command, args []string) error {
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("load config: %w", err)
	}

	if err := ensureSetup(cfg); err != nil {
		return err
	}

	provider, err := llm.NewProvider(cfg)
	if err != nil {
		return err
	}

	// Build the prompt from args, stdin, and file flag
	var parts []string

	// Check stdin
	stat, _ := os.Stdin.Stat()
	if (stat.Mode() & os.ModeCharDevice) == 0 {
		data, err := io.ReadAll(os.Stdin)
		if err != nil {
			return fmt.Errorf("read stdin: %w", err)
		}
		if len(data) > 0 {
			parts = append(parts, "Context from stdin:\n```\n"+strings.TrimSpace(string(data))+"\n```")
		}
	}

	// Check file flag
	filePath, _ := cmd.Flags().GetString("file")
	if filePath != "" {
		data, err := os.ReadFile(filePath)
		if err != nil {
			return fmt.Errorf("read file %s: %w", filePath, err)
		}
		parts = append(parts, fmt.Sprintf("Contents of %s:\n```\n%s\n```", filePath, strings.TrimSpace(string(data))))
	}

	// Question from args
	if len(args) > 0 {
		parts = append(parts, strings.Join(args, " "))
	}

	if len(parts) == 0 {
		return fmt.Errorf("no question provided. Usage: bro ask \"your question\"")
	}

	prompt := strings.Join(parts, "\n\n")
	systemPrompt := "You are a helpful developer assistant. Be concise and precise."

	ctx := context.Background()
	ch, err := provider.StreamGenerate(ctx, prompt, systemPrompt)
	if err != nil {
		return fmt.Errorf("generate response: %w", err)
	}

	for chunk := range ch {
		fmt.Print(chunk)
	}
	fmt.Println()

	return nil
}
