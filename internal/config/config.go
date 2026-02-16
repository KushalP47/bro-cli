package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/viper"
)

// Config holds the application configuration.
type Config struct {
	Provider    string    `yaml:"provider" mapstructure:"provider"`
	Anthropic   APIConfig `yaml:"anthropic" mapstructure:"anthropic"`
	OpenAI      APIConfig `yaml:"openai" mapstructure:"openai"`
	Profiles    []Profile `yaml:"profiles" mapstructure:"profiles"`
	StandupRepos []string `yaml:"standup_repos" mapstructure:"standup_repos"`
}

// APIConfig holds API credentials for an LLM provider.
type APIConfig struct {
	APIKey string `yaml:"api_key" mapstructure:"api_key"`
	Model  string `yaml:"model" mapstructure:"model"`
}

// Load reads the config from ~/.bro/config.yaml.
func Load() (*Config, error) {
	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("unmarshal config: %w", err)
	}
	return &cfg, nil
}

// Dir returns the path to the bro config directory (~/.bro).
func Dir() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("find home directory: %w", err)
	}
	return filepath.Join(home, ".bro"), nil
}

// EnsureDir creates the config directory and subdirectories if they don't exist.
func EnsureDir() error {
	dir, err := Dir()
	if err != nil {
		return err
	}
	dirs := []string{
		dir,
		filepath.Join(dir, "standups"),
		filepath.Join(dir, "templates"),
	}
	for _, d := range dirs {
		if err := os.MkdirAll(d, 0755); err != nil {
			return fmt.Errorf("create directory %s: %w", d, err)
		}
	}
	return nil
}
