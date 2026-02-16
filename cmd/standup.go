package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var standupCmd = &cobra.Command{
	Use:   "standup",
	Short: "Generate a standup summary from git logs",
	Long:  "Scans configured repos for recent git activity and generates an LLM-summarized standup report.",
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("bro standup — not yet implemented")
		return nil
	},
}

func init() {
	rootCmd.AddCommand(standupCmd)
}
