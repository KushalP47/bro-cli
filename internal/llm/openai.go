package llm

import "context"

// OpenAIProvider implements the Provider interface using the OpenAI Chat Completions API.
type OpenAIProvider struct {
	APIKey string
	Model  string
}

func NewOpenAIProvider(apiKey, model string) *OpenAIProvider {
	if model == "" {
		model = "gpt-4o"
	}
	return &OpenAIProvider{APIKey: apiKey, Model: model}
}

func (o *OpenAIProvider) Generate(ctx context.Context, prompt, systemPrompt string) (string, error) {
	// TODO: implement OpenAI Chat Completions API call
	return "", nil
}

func (o *OpenAIProvider) StreamGenerate(ctx context.Context, prompt, systemPrompt string) (<-chan string, error) {
	// TODO: implement streaming OpenAI Chat Completions API call
	return nil, nil
}
