// Package ai contains the Gemini-backed RAG / AI Data Assistant building blocks.
//
// gemini_client.go is a thin REST client over Google's Generative Language API
// (https://generativelanguage.googleapis.com/v1beta). We use net/http directly
// instead of an SDK so the binary stays small (Lambda-friendly) and to keep
// transitive dependencies minimal.
package ai

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

// ErrGeminiUnavailable is returned when the Gemini API is unreachable or
// returns a non-recoverable error after retries. Higher layers should map this
// to a graceful 503 to the user.
var ErrGeminiUnavailable = errors.New("gemini api unavailable")

const geminiBaseURL = "https://generativelanguage.googleapis.com/v1beta"

// GeminiClient performs HTTP calls against the Gemini REST API.
type GeminiClient struct {
	apiKey     string
	chatModel  string
	embedModel string
	http       *http.Client
}

// NewGeminiClient constructs a GeminiClient with the configured timeout and models.
func NewGeminiClient(apiKey, chatModel, embedModel string, timeout time.Duration) *GeminiClient {
	if chatModel == "" {
		chatModel = "gemini-1.5-flash"
	}
	if embedModel == "" {
		embedModel = "text-embedding-004"
	}
	if timeout <= 0 {
		timeout = 25 * time.Second
	}
	return &GeminiClient{
		apiKey:     apiKey,
		chatModel:  chatModel,
		embedModel: embedModel,
		http:       &http.Client{Timeout: timeout},
	}
}

// Configured reports whether the client has an API key. When false, callers
// should skip AI work (no-op) instead of returning an error to the user.
func (c *GeminiClient) Configured() bool { return c != nil && strings.TrimSpace(c.apiKey) != "" }

// ChatModel returns the currently configured chat model name.
func (c *GeminiClient) ChatModel() string { return c.chatModel }

// EmbedModel returns the currently configured embedding model name.
func (c *GeminiClient) EmbedModel() string { return c.embedModel }

// ----------------------------------------------------------------------------
// Embeddings
// ----------------------------------------------------------------------------

type embedRequest struct {
	Model   string         `json:"model"`
	Content embedRequestIn `json:"content"`
}

type embedRequestIn struct {
	Parts []geminiPart `json:"parts"`
}

type embedResponse struct {
	Embedding struct {
		Values []float32 `json:"values"`
	} `json:"embedding"`
}

// Embed returns the 768-dim embedding for the given text using text-embedding-004.
// Empty input returns an error so callers don't silently store a zero vector.
func (c *GeminiClient) Embed(ctx context.Context, text string) ([]float32, error) {
	text = strings.TrimSpace(text)
	if text == "" {
		return nil, errors.New("embed: empty text")
	}
	if !c.Configured() {
		return nil, ErrGeminiUnavailable
	}

	body := embedRequest{
		Model: "models/" + c.embedModel,
		Content: embedRequestIn{
			Parts: []geminiPart{{Text: text}},
		},
	}
	url := fmt.Sprintf("%s/models/%s:embedContent?key=%s", geminiBaseURL, c.embedModel, c.apiKey)

	var out embedResponse
	if err := c.doJSON(ctx, http.MethodPost, url, body, &out); err != nil {
		return nil, err
	}
	if len(out.Embedding.Values) == 0 {
		return nil, fmt.Errorf("embed: empty response")
	}
	return out.Embedding.Values, nil
}

// ----------------------------------------------------------------------------
// Chat / generateContent (with function calling support)
// ----------------------------------------------------------------------------

// GeminiPart represents a single part of a content item. Exactly one of Text,
// FunctionCall, or FunctionResponse should be non-nil per part.
type geminiPart struct {
	Text             string                  `json:"text,omitempty"`
	FunctionCall     *GeminiFunctionCall     `json:"functionCall,omitempty"`
	FunctionResponse *GeminiFunctionResponse `json:"functionResponse,omitempty"`
}

// GeminiContent is one turn in the conversation (user / model / function).
type GeminiContent struct {
	Role  string       `json:"role"`
	Parts []geminiPart `json:"parts"`
}

// GeminiFunctionCall mirrors the model's request to invoke a tool.
type GeminiFunctionCall struct {
	Name string                 `json:"name"`
	Args map[string]interface{} `json:"args"`
}

// GeminiFunctionResponse is what we send back after running a tool.
type GeminiFunctionResponse struct {
	Name     string                 `json:"name"`
	Response map[string]interface{} `json:"response"`
}

// GeminiTool is a top-level tool declaration container.
type GeminiTool struct {
	FunctionDeclarations []GeminiFunctionDeclaration `json:"functionDeclarations"`
}

// GeminiFunctionDeclaration declares a single callable tool to the model.
type GeminiFunctionDeclaration struct {
	Name        string         `json:"name"`
	Description string         `json:"description"`
	Parameters  GeminiSchema   `json:"parameters"`
}

// GeminiSchema is a (subset of) OpenAPI schema accepted by Gemini for tool
// parameter definitions.
type GeminiSchema struct {
	Type        string                  `json:"type,omitempty"`
	Description string                  `json:"description,omitempty"`
	Enum        []string                `json:"enum,omitempty"`
	Properties  map[string]GeminiSchema `json:"properties,omitempty"`
	Required    []string                `json:"required,omitempty"`
	Items       *GeminiSchema           `json:"items,omitempty"`
}

type geminiSystemInstruction struct {
	Parts []geminiPart `json:"parts"`
}

// GeminiRequest is the body sent to :generateContent.
type GeminiRequest struct {
	SystemInstruction *geminiSystemInstruction `json:"systemInstruction,omitempty"`
	Contents          []GeminiContent          `json:"contents"`
	Tools             []GeminiTool             `json:"tools,omitempty"`
	GenerationConfig  *GeminiGenerationConfig  `json:"generationConfig,omitempty"`
}

// GeminiGenerationConfig is a small subset of the upstream config fields.
type GeminiGenerationConfig struct {
	Temperature     *float64 `json:"temperature,omitempty"`
	MaxOutputTokens *int     `json:"maxOutputTokens,omitempty"`
}

// GeminiCandidate is one generation option.
type GeminiCandidate struct {
	Content      GeminiContent `json:"content"`
	FinishReason string        `json:"finishReason,omitempty"`
}

// GeminiResponse is the parsed :generateContent response.
type GeminiResponse struct {
	Candidates []GeminiCandidate `json:"candidates"`
}

// FirstCandidate returns the first candidate or nil.
func (r *GeminiResponse) FirstCandidate() *GeminiCandidate {
	if len(r.Candidates) == 0 {
		return nil
	}
	return &r.Candidates[0]
}

// FunctionCalls extracts every functionCall part from the first candidate.
func (r *GeminiResponse) FunctionCalls() []GeminiFunctionCall {
	c := r.FirstCandidate()
	if c == nil {
		return nil
	}
	var out []GeminiFunctionCall
	for _, p := range c.Content.Parts {
		if p.FunctionCall != nil {
			out = append(out, *p.FunctionCall)
		}
	}
	return out
}

// Text concatenates every text part of the first candidate.
func (r *GeminiResponse) Text() string {
	c := r.FirstCandidate()
	if c == nil {
		return ""
	}
	var b strings.Builder
	for _, p := range c.Content.Parts {
		if p.Text != "" {
			if b.Len() > 0 {
				b.WriteString("\n")
			}
			b.WriteString(p.Text)
		}
	}
	return b.String()
}

// NewSystemInstruction wraps a plain text system prompt.
func NewSystemInstruction(prompt string) *geminiSystemInstruction {
	return &geminiSystemInstruction{Parts: []geminiPart{{Text: prompt}}}
}

// NewUserContent builds a user-role content item with a single text part.
func NewUserContent(text string) GeminiContent {
	return GeminiContent{Role: "user", Parts: []geminiPart{{Text: text}}}
}

// NewModelContent builds a model-role content item with a single text part.
func NewModelContent(text string) GeminiContent {
	return GeminiContent{Role: "model", Parts: []geminiPart{{Text: text}}}
}

// NewFunctionResponseContent encodes our tool-call result back to the model.
// Gemini expects role "function" for function responses.
func NewFunctionResponseContent(name string, response map[string]interface{}) GeminiContent {
	return GeminiContent{
		Role: "function",
		Parts: []geminiPart{{
			FunctionResponse: &GeminiFunctionResponse{Name: name, Response: response},
		}},
	}
}

// GenerateContent calls models/{chat}:generateContent. Caller owns history.
func (c *GeminiClient) GenerateContent(ctx context.Context, req GeminiRequest) (*GeminiResponse, error) {
	if !c.Configured() {
		return nil, ErrGeminiUnavailable
	}
	url := fmt.Sprintf("%s/models/%s:generateContent?key=%s", geminiBaseURL, c.chatModel, c.apiKey)
	var out GeminiResponse
	if err := c.doJSON(ctx, http.MethodPost, url, req, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// ----------------------------------------------------------------------------
// HTTP plumbing with one retry on 429/5xx
// ----------------------------------------------------------------------------

func (c *GeminiClient) doJSON(ctx context.Context, method, url string, body any, out any) error {
	const maxAttempts = 2
	var lastErr error
	for attempt := 1; attempt <= maxAttempts; attempt++ {
		err := c.doOnce(ctx, method, url, body, out)
		if err == nil {
			return nil
		}
		lastErr = err
		// Only retry on transient errors signalled by the wrapped error message.
		if !isTransient(err) || attempt == maxAttempts {
			break
		}
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(time.Duration(attempt) * 400 * time.Millisecond):
		}
	}
	return lastErr
}

func (c *GeminiClient) doOnce(ctx context.Context, method, url string, body any, out any) error {
	var reqBody io.Reader
	if body != nil {
		raw, err := json.Marshal(body)
		if err != nil {
			return fmt.Errorf("gemini: marshal request: %w", err)
		}
		reqBody = bytes.NewReader(raw)
	}
	req, err := http.NewRequestWithContext(ctx, method, url, reqBody)
	if err != nil {
		return fmt.Errorf("gemini: build request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	resp, err := c.http.Do(req)
	if err != nil {
		return fmt.Errorf("%w: %v", ErrGeminiUnavailable, err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("gemini: read body: %w", err)
	}

	if resp.StatusCode >= 400 {
		// 429 or 5xx are considered transient.
		if resp.StatusCode == http.StatusTooManyRequests || resp.StatusCode >= 500 {
			return fmt.Errorf("%w: gemini %d: %s", ErrGeminiUnavailable, resp.StatusCode, truncate(string(respBody), 300))
		}
		return fmt.Errorf("gemini: %d: %s", resp.StatusCode, truncate(string(respBody), 300))
	}
	if out == nil {
		return nil
	}
	if err := json.Unmarshal(respBody, out); err != nil {
		return fmt.Errorf("gemini: decode body: %w", err)
	}
	return nil
}

func isTransient(err error) bool {
	if err == nil {
		return false
	}
	return errors.Is(err, ErrGeminiUnavailable)
}

func truncate(s string, n int) string {
	if len(s) <= n {
		return s
	}
	return s[:n] + "..."
}
