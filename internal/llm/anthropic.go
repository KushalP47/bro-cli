package llm

import "context"

// AnthropicProvider implements the Provider interface using the Anthropic Messages API.
type AnthropicProvider struct {
	APIKey string
	Model  string
}

func NewAnthropicProvider(apiKey, model string) *AnthropicProvider {
	if model == "" {
		model = "claude-sonnet-4-5-20250929"
	}
	return &AnthropicProvider{APIKey: apiKey, Model: model}
}

func (a *AnthropicProvider) Generate(ctx context.Context, prompt, systemPrompt string) (string, error) {
	// TODO: implement Anthropic Messages API call
	return "", nil
}

func (a *AnthropicProvider) StreamGenerate(ctx context.Context, prompt, systemPrompt string) (<-chan string, error) {
	// TODO: implement streaming Anthropic Messages API call
	return nil, nil
}
