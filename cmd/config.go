package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Manage bro configuration",
	Long:  "View and modify config settings, profiles, and API keys.",
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("bro config — not yet implemented")
		return nil
	},
}

func init() {
	rootCmd.AddCommand(configCmd)
}
