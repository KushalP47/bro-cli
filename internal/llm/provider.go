package llm

import (
	"context"
	"fmt"
	"os"

	"github.com/kushalpatel/bro-cli/internal/config"
)

// Provider is the interface for LLM backends.
type Provider interface {
	Generate(ctx context.Context, prompt, systemPrompt string) (string, error)
	StreamGenerate(ctx context.Context, prompt, systemPrompt string) (<-chan string, error)
}

// NewProvider creates an LLM provider based on the config.
func NewProvider(cfg *config.Config) (Provider, error) {
	switch cfg.Provider {
	case "anthropic", "":
		apiKey := cfg.Anthropic.APIKey
		if apiKey == "" {
			apiKey = os.Getenv("BRO_ANTHROPIC_API_KEY")
		}
		if apiKey == "" {
			return nil, fmt.Errorf("no Anthropic API key configured. Run 'bro config init' or set BRO_ANTHROPIC_API_KEY")
		}
		return NewAnthropicProvider(apiKey, cfg.Anthropic.Model), nil

	case "openai":
		apiKey := cfg.OpenAI.APIKey
		if apiKey == "" {
			apiKey = os.Getenv("BRO_OPENAI_API_KEY")
		}
		if apiKey == "" {
			return nil, fmt.Errorf("no OpenAI API key configured. Run 'bro config init' or set BRO_OPENAI_API_KEY")
		}
		return NewOpenAIProvider(apiKey, cfg.OpenAI.Model), nil

	default:
		return nil, fmt.Errorf("unknown LLM provider: %s", cfg.Provider)
	}
}
