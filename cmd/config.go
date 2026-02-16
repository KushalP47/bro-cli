package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/kushalpatel/bro-cli/internal/config"
	"github.com/kushalpatel/bro-cli/internal/ui"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Manage bro configuration",
	Long:  "View and modify config settings, profiles, and API keys.",
}

var configShowCmd = &cobra.Command{
	Use:   "show",
	Short: "Show current configuration",
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := config.Load()
		if err != nil {
			return fmt.Errorf("load config: %w", err)
		}

		fmt.Printf("Provider: %s\n", defaultStr(cfg.Provider, "anthropic"))
		if cfg.Anthropic.APIKey != "" {
			fmt.Printf("Anthropic API Key: %s...%s\n", cfg.Anthropic.APIKey[:4], cfg.Anthropic.APIKey[len(cfg.Anthropic.APIKey)-4:])
		}
		if cfg.Anthropic.Model != "" {
			fmt.Printf("Anthropic Model: %s\n", cfg.Anthropic.Model)
		}
		if cfg.OpenAI.APIKey != "" {
			fmt.Printf("OpenAI API Key: %s...%s\n", cfg.OpenAI.APIKey[:4], cfg.OpenAI.APIKey[len(cfg.OpenAI.APIKey)-4:])
		}
		if cfg.OpenAI.Model != "" {
			fmt.Printf("OpenAI Model: %s\n", cfg.OpenAI.Model)
		}
		if len(cfg.Profiles) > 0 {
			fmt.Println("\nProfiles:")
			for _, p := range cfg.Profiles {
				fmt.Printf("  - %s (%s) → %s\n", p.Name, p.CommitStyle, p.DirPattern)
			}
		}
		if len(cfg.StandupRepos) > 0 {
			fmt.Println("\nStandup Repos:")
			for _, r := range cfg.StandupRepos {
				fmt.Printf("  - %s\n", r)
			}
		}
		return nil
	},
}

var configSetCmd = &cobra.Command{
	Use:   "set <key> <value>",
	Short: "Set a config value",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		key, value := args[0], args[1]
		viper.Set(key, value)

		if err := config.EnsureDir(); err != nil {
			return err
		}

		configDir, err := config.Dir()
		if err != nil {
			return err
		}

		if err := viper.WriteConfigAs(configDir + "/config.yaml"); err != nil {
			return fmt.Errorf("write config: %w", err)
		}

		fmt.Printf("Set %s = %s\n", key, value)
		return nil
	},
}

var configInitCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize configuration with setup wizard",
	RunE: func(cmd *cobra.Command, args []string) error {
		return runSetupWizard()
	},
}

func init() {
	configCmd.AddCommand(configShowCmd)
	configCmd.AddCommand(configSetCmd)
	configCmd.AddCommand(configInitCmd)
	rootCmd.AddCommand(configCmd)
}

func runSetupWizard() error {
	reader := bufio.NewReader(os.Stdin)

	fmt.Println()
	fmt.Println(ui.TitleStyle.Render("Welcome to bro! Let's get you set up."))
	fmt.Println()

	// Choose provider
	fmt.Print("LLM Provider [anthropic/openai] (default: anthropic): ")
	provider, _ := reader.ReadString('\n')
	provider = strings.TrimSpace(provider)
	if provider == "" {
		provider = "anthropic"
	}
	if provider != "anthropic" && provider != "openai" {
		return fmt.Errorf("invalid provider: %s. Choose 'anthropic' or 'openai'", provider)
	}

	// Get API key
	fmt.Printf("API Key for %s: ", provider)
	apiKey, _ := reader.ReadString('\n')
	apiKey = strings.TrimSpace(apiKey)
	if apiKey == "" {
		return fmt.Errorf("API key cannot be empty")
	}

	// Get model (optional)
	var defaultModel string
	if provider == "anthropic" {
		defaultModel = "claude-sonnet-4-5-20250929"
	} else {
		defaultModel = "gpt-4o"
	}
	fmt.Printf("Model (default: %s): ", defaultModel)
	model, _ := reader.ReadString('\n')
	model = strings.TrimSpace(model)
	if model == "" {
		model = defaultModel
	}

	// Save config
	viper.Set("provider", provider)
	if provider == "anthropic" {
		viper.Set("anthropic.api_key", apiKey)
		viper.Set("anthropic.model", model)
	} else {
		viper.Set("openai.api_key", apiKey)
		viper.Set("openai.model", model)
	}

	if err := config.EnsureDir(); err != nil {
		return err
	}

	configDir, err := config.Dir()
	if err != nil {
		return err
	}

	if err := viper.WriteConfigAs(configDir + "/config.yaml"); err != nil {
		return fmt.Errorf("write config: %w", err)
	}

	fmt.Println()
	fmt.Println(ui.SuccessStyle.Render("Configuration saved!"))
	fmt.Printf("Config file: %s/config.yaml\n", configDir)

	return nil
}

func defaultStr(s, def string) string {
	if s == "" {
		return def
	}
	return s
}
