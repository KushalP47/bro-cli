package llm

import "context"

// Provider is the interface for LLM backends.
type Provider interface {
	Generate(ctx context.Context, prompt, systemPrompt string) (string, error)
	StreamGenerate(ctx context.Context, prompt, systemPrompt string) (<-chan string, error)
}
