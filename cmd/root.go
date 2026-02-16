package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/kushalpatel/bro-cli/internal/config"
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

// ensureSetup checks if config exists and runs setup wizard on first run.
func ensureSetup(cfg *config.Config) error {
	configDir, err := config.Dir()
	if err != nil {
		return err
	}

	configPath := filepath.Join(configDir, "config.yaml")
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		// No config file — check if API key is available via env
		if os.Getenv("BRO_ANTHROPIC_API_KEY") != "" || os.Getenv("BRO_OPENAI_API_KEY") != "" {
			return nil
		}

		fmt.Println("No configuration found. Starting setup wizard...")
		fmt.Println()
		return runSetupWizard()
	}

	return nil
}
