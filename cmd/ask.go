package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var askCmd = &cobra.Command{
	Use:   "ask [question]",
	Short: "Ask the LLM a question",
	Long:  "General LLM queries with support for stdin pipe and -f file flag. Streams response to terminal.",
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("bro ask — not yet implemented")
		return nil
	},
}

func init() {
	askCmd.Flags().StringP("file", "f", "", "File to include as context")
	rootCmd.AddCommand(askCmd)
}
