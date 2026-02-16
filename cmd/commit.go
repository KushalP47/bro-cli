package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var commitCmd = &cobra.Command{
	Use:   "commit",
	Short: "Generate a commit message from staged changes",
	Long:  "Analyzes git diff --staged and uses an LLM to generate a commit message. Supports accept, edit, and regenerate flow.",
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("bro commit — not yet implemented")
		return nil
	},
}

func init() {
	rootCmd.AddCommand(commitCmd)
}
