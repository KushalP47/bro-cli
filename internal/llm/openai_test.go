package llm

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestOpenAIGenerate(t *testing.T) {
	tests := []struct {
		name     string
		status   int
		response any
		wantErr  bool
		wantText string
	}{
		{
			name:   "successful response",
			status: http.StatusOK,
			response: openaiResponse{
				Choices: []struct {
					Message struct {
						Content string `json:"content"`
					} `json:"message"`
				}{
					{Message: struct {
						Content string `json:"content"`
					}{Content: "Hello from OpenAI!"}},
				},
			},
			wantText: "Hello from OpenAI!",
		},
		{
			name:   "empty choices",
			status: http.StatusOK,
			response: openaiResponse{
				Choices: []struct {
					Message struct {
						Content string `json:"content"`
					} `json:"message"`
				}{},
			},
			wantErr: true,
		},
		{
			name:     "API error",
			status:   http.StatusUnauthorized,
			response: map[string]any{"error": map[string]string{"message": "invalid api key"}},
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.Header.Get("Authorization") != "Bearer test-key" {
					t.Error("expected Bearer authorization header")
				}

				var req openaiRequest
				json.NewDecoder(r.Body).Decode(&req)
				if req.Model != "test-model" {
					t.Errorf("expected model test-model, got %s", req.Model)
				}

				w.WriteHeader(tt.status)
				json.NewEncoder(w).Encode(tt.response)
			}))
			defer server.Close()

			provider := &OpenAIProvider{APIKey: "test-key", Model: "test-model"}

			origTransport := http.DefaultTransport
			http.DefaultTransport = &rewriteTransport{
				host:      strings.TrimPrefix(server.URL, "http://"),
				transport: origTransport,
			}
			defer func() { http.DefaultTransport = origTransport }()

			result, err := provider.Generate(context.Background(), "test prompt", "test system")
			if tt.wantErr {
				if err == nil {
					t.Error("expected error, got nil")
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if result != tt.wantText {
				t.Errorf("expected %q, got %q", tt.wantText, result)
			}
		})
	}
}

func TestOpenAIStreamGenerate(t *testing.T) {
	sseData := `data: {"choices":[{"delta":{"content":"Hello"}}]}

data: {"choices":[{"delta":{"content":" world"}}]}

data: [DONE]

`
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/event-stream")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(sseData))
	}))
	defer server.Close()

	provider := &OpenAIProvider{APIKey: "test-key", Model: "test-model"}

	origTransport := http.DefaultTransport
	http.DefaultTransport = &rewriteTransport{
		host:      strings.TrimPrefix(server.URL, "http://"),
		transport: origTransport,
	}
	defer func() { http.DefaultTransport = origTransport }()

	ch, err := provider.StreamGenerate(context.Background(), "test", "system")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	var chunks []string
	for chunk := range ch {
		chunks = append(chunks, chunk)
	}

	if len(chunks) != 2 {
		t.Fatalf("expected 2 chunks, got %d: %v", len(chunks), chunks)
	}
	if chunks[0] != "Hello" || chunks[1] != " world" {
		t.Errorf("unexpected chunks: %v", chunks)
	}
}
