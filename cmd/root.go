package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var rootCmd = &cobra.Command{
	Use:   "bro",
	Short: "Developer productivity CLI",
	Long:  "bro — smart commit messages, command shortcuts, standup summaries, and AI-assisted Q&A.",
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)
}

func initConfig() {
	home, err := os.UserHomeDir()
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error finding home directory:", err)
		os.Exit(1)
	}

	viper.AddConfigPath(home + "/.bro")
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")

	viper.SetEnvPrefix("BRO")
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		// Config file not found is fine on first run
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			fmt.Fprintln(os.Stderr, "Error reading config:", err)
		}
	}
}
