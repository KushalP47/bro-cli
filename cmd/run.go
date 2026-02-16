package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var runCmd = &cobra.Command{
	Use:   "run [alias]",
	Short: "Run a command shortcut",
	Long:  "Fuzzy-searchable command shortcuts with arrow key navigation, or direct alias execution.",
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("bro run — not yet implemented")
		return nil
	},
}

func init() {
	runCmd.Flags().String("preset", "", "Preset parameters for the command")
	rootCmd.AddCommand(runCmd)
}
