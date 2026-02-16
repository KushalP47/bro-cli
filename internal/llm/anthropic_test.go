package llm

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestAnthropicGenerate(t *testing.T) {
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
			response: anthropicResponse{
				Content: []struct {
					Text string `json:"text"`
				}{
					{Text: "Hello, world!"},
				},
			},
			wantText: "Hello, world!",
		},
		{
			name:   "empty content",
			status: http.StatusOK,
			response: anthropicResponse{
				Content: []struct {
					Text string `json:"text"`
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
				if r.Header.Get("x-api-key") != "test-key" {
					t.Error("expected x-api-key header")
				}
				if r.Header.Get("anthropic-version") != "2023-06-01" {
					t.Error("expected anthropic-version header")
				}

				var req anthropicRequest
				json.NewDecoder(r.Body).Decode(&req)
				if req.Model != "test-model" {
					t.Errorf("expected model test-model, got %s", req.Model)
				}
				if len(req.Messages) != 1 || req.Messages[0].Role != "user" {
					t.Error("expected single user message")
				}

				w.WriteHeader(tt.status)
				json.NewEncoder(w).Encode(tt.response)
			}))
			defer server.Close()

			provider := &AnthropicProvider{APIKey: "test-key", Model: "test-model"}

			origTransport := http.DefaultTransport
			http.DefaultTransport = &rewriteTransport{
				host:      strings.TrimPrefix(server.URL, "http://"),
				transport: http.DefaultTransport,
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

func TestAnthropicStreamGenerate(t *testing.T) {
	sseData := `event: content_block_start
data: {"type":"content_block_start","index":0}

event: content_block_delta
data: {"type":"content_block_delta","delta":{"type":"text_delta","text":"Hello"}}

event: content_block_delta
data: {"type":"content_block_delta","delta":{"type":"text_delta","text":" world"}}

event: message_stop
data: {"type":"message_stop"}

`
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/event-stream")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(sseData))
	}))
	defer server.Close()

	provider := &AnthropicProvider{APIKey: "test-key", Model: "test-model"}

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

// rewriteTransport redirects requests to the test server
type rewriteTransport struct {
	host      string
	transport http.RoundTripper
}

func (t *rewriteTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	req.URL.Scheme = "http"
	req.URL.Host = t.host
	return t.transport.RoundTrip(req)
}
