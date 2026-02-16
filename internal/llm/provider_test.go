package llm

import (
	"os"
	"testing"

	"github.com/kushalpatel/bro-cli/internal/config"
)

func TestNewProvider(t *testing.T) {
	tests := []struct {
		name     string
		cfg      *config.Config
		envKey   string
		envValue string
		wantType string
		wantErr  bool
	}{
		{
			name: "anthropic with config key",
			cfg: &config.Config{
				Provider:  "anthropic",
				Anthropic: config.APIConfig{APIKey: "sk-test", Model: "claude-sonnet-4-5-20250929"},
			},
			wantType: "*llm.AnthropicProvider",
		},
		{
			name: "openai with config key",
			cfg: &config.Config{
				Provider: "openai",
				OpenAI:   config.APIConfig{APIKey: "sk-test", Model: "gpt-4o"},
			},
			wantType: "*llm.OpenAIProvider",
		},
		{
			name:     "anthropic with env key",
			cfg:      &config.Config{Provider: "anthropic"},
			envKey:   "BRO_ANTHROPIC_API_KEY",
			envValue: "sk-env-test",
			wantType: "*llm.AnthropicProvider",
		},
		{
			name:    "missing anthropic key",
			cfg:     &config.Config{Provider: "anthropic"},
			wantErr: true,
		},
		{
			name:    "missing openai key",
			cfg:     &config.Config{Provider: "openai"},
			wantErr: true,
		},
		{
			name:    "unknown provider",
			cfg:     &config.Config{Provider: "gemini"},
			wantErr: true,
		},
		{
			name: "empty provider defaults to anthropic",
			cfg: &config.Config{
				Provider:  "",
				Anthropic: config.APIConfig{APIKey: "sk-test"},
			},
			wantType: "*llm.AnthropicProvider",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.envKey != "" {
				os.Setenv(tt.envKey, tt.envValue)
				defer os.Unsetenv(tt.envKey)
			}

			provider, err := NewProvider(tt.cfg)
			if tt.wantErr {
				if err == nil {
					t.Error("expected error, got nil")
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if provider == nil {
				t.Fatal("expected non-nil provider")
			}
		})
	}
}
