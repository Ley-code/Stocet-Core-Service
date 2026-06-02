// Package ai — claude_client.go is a thin REST client over Anthropic's Messages API
// (https://api.anthropic.com/v1/messages). We use net/http directly (no SDK) to keep
// the binary Lambda-friendly. Embeddings remain on GeminiClient.
package ai

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"time"
)

// ErrClaudeUnavailable is returned when the Anthropic API is unreachable or returns
// a non-recoverable error after retries.
var ErrClaudeUnavailable = errors.New("claude api unavailable")

const claudeBaseURL = "https://api.anthropic.com/v1"
const claudeAPIVersion = "2023-06-01"

// claudeDefaultMaxTokens caps output per response; Haiku 4.5 supports up to 64K.
const claudeDefaultMaxTokens = 8192

// ClaudeClient performs HTTP calls against the Anthropic Messages API.
type ClaudeClient struct {
	apiKey    string
	chatModel string
	http      *http.Client
}

// NewClaudeClient constructs a ClaudeClient. model defaults to "claude-haiku-4-5".
func NewClaudeClient(apiKey, model string, timeout time.Duration) *ClaudeClient {
	if model == "" {
		model = "claude-haiku-4-5"
	}
	if timeout <= 0 {
		timeout = 25 * time.Second
	}
	return &ClaudeClient{
		apiKey:    apiKey,
		chatModel: model,
		http:      &http.Client{Timeout: timeout},
	}
}

// Configured reports whether the client has an API key set.
func (c *ClaudeClient) Configured() bool { return c != nil && strings.TrimSpace(c.apiKey) != "" }

// ChatModel returns the configured model name.
func (c *ClaudeClient) ChatModel() string { return c.chatModel }

// ─────────────────────────────────────────────────────────────────────────────
// Request / Response types
// ─────────────────────────────────────────────────────────────────────────────

// ClaudeCacheControl marks a content block for prompt caching.
type ClaudeCacheControl struct {
	Type string `json:"type"` // "ephemeral"
}

// ClaudeSystemBlock is one entry in the top-level "system" array.
// Setting CacheControl enables prompt caching for the system prompt.
type ClaudeSystemBlock struct {
	Type         string              `json:"type"`
	Text         string              `json:"text"`
	CacheControl *ClaudeCacheControl `json:"cache_control,omitempty"`
}

// ClaudeContentBlock represents any block inside a message's content array:
// text, tool_use (from model), or tool_result (from us).
type ClaudeContentBlock struct {
	Type      string                 `json:"type"`
	Text      string                 `json:"text,omitempty"`
	ID        string                 `json:"id,omitempty"`        // tool_use id
	Name      string                 `json:"name,omitempty"`      // tool_use name
	Input     map[string]interface{} `json:"input,omitempty"`     // tool_use input
	ToolUseID string                 `json:"tool_use_id,omitempty"` // tool_result back-reference
	Content   interface{}            `json:"content,omitempty"`   // tool_result payload (string)
}

// ClaudeMessage is one turn in the conversation. Content is either a plain string
// (for simple user messages) or []ClaudeContentBlock (for tool_use / tool_result turns).
type ClaudeMessage struct {
	Role    string      `json:"role"`
	Content interface{} `json:"content"`
}

// ClaudeSchema is a subset of JSON Schema accepted by Anthropic for tool input_schema.
type ClaudeSchema struct {
	Type        string                  `json:"type,omitempty"`
	Description string                  `json:"description,omitempty"`
	Enum        []string                `json:"enum,omitempty"`
	Properties  map[string]ClaudeSchema `json:"properties,omitempty"`
	Required    []string                `json:"required,omitempty"`
	Items       *ClaudeSchema           `json:"items,omitempty"`
}

// ClaudeTool declares one callable tool to the model.
type ClaudeTool struct {
	Name        string       `json:"name"`
	Description string       `json:"description"`
	InputSchema ClaudeSchema `json:"input_schema"`
}

// ClaudeRequest is the body sent to POST /v1/messages.
type ClaudeRequest struct {
	Model     string              `json:"model"`
	MaxTokens int                 `json:"max_tokens"`
	System    []ClaudeSystemBlock `json:"system,omitempty"`
	Messages  []ClaudeMessage     `json:"messages"`
	Tools     []ClaudeTool        `json:"tools,omitempty"`
}

// ClaudeToolUse is a single tool invocation extracted from the model's response.
type ClaudeToolUse struct {
	ID    string                 `json:"id"`
	Name  string                 `json:"name"`
	Input map[string]interface{} `json:"input"`
}

// ClaudeResponse is the parsed /v1/messages response.
type ClaudeResponse struct {
	ID         string               `json:"id"`
	Type       string               `json:"type"`
	Role       string               `json:"role"`
	Content    []ClaudeContentBlock `json:"content"`
	Model      string               `json:"model"`
	StopReason string               `json:"stop_reason"`
	Usage      struct {
		InputTokens              int `json:"input_tokens"`
		OutputTokens             int `json:"output_tokens"`
		CacheCreationInputTokens int `json:"cache_creation_input_tokens"`
		CacheReadInputTokens     int `json:"cache_read_input_tokens"`
	} `json:"usage"`
}

// Text concatenates all text blocks from the response.
func (r *ClaudeResponse) Text() string {
	if r == nil {
		return ""
	}
	var b strings.Builder
	for _, block := range r.Content {
		if block.Type == "text" && block.Text != "" {
			if b.Len() > 0 {
				b.WriteString("\n")
			}
			b.WriteString(block.Text)
		}
	}
	return b.String()
}

// ToolUses extracts every tool_use block from the response.
func (r *ClaudeResponse) ToolUses() []ClaudeToolUse {
	if r == nil {
		return nil
	}
	var out []ClaudeToolUse
	for _, block := range r.Content {
		if block.Type == "tool_use" {
			out = append(out, ClaudeToolUse{
				ID:    block.ID,
				Name:  block.Name,
				Input: block.Input,
			})
		}
	}
	return out
}

// ToAssistantMessage converts the response content into an assistant turn for
// echoing back in the next request (required by the tool-calling loop).
func (r *ClaudeResponse) ToAssistantMessage() ClaudeMessage {
	return ClaudeMessage{Role: "assistant", Content: r.Content}
}

// ─────────────────────────────────────────────────────────────────────────────
// Helper constructors
// ─────────────────────────────────────────────────────────────────────────────

// NewClaudeUserMessage builds a plain text user turn.
func NewClaudeUserMessage(text string) ClaudeMessage {
	return ClaudeMessage{Role: "user", Content: text}
}

// NewClaudeAssistantMessage builds a plain text assistant turn (for chat history).
func NewClaudeAssistantMessage(text string) ClaudeMessage {
	return ClaudeMessage{
		Role:    "assistant",
		Content: []ClaudeContentBlock{{Type: "text", Text: text}},
	}
}

// NewClaudeSystemBlock returns a system block with prompt caching enabled.
func NewClaudeSystemBlock(text string) ClaudeSystemBlock {
	return ClaudeSystemBlock{
		Type:         "text",
		Text:         text,
		CacheControl: &ClaudeCacheControl{Type: "ephemeral"},
	}
}

// NewClaudeToolResultMessage wraps tool result blocks in a user turn.
func NewClaudeToolResultMessage(results []ClaudeContentBlock) ClaudeMessage {
	return ClaudeMessage{Role: "user", Content: results}
}

// ─────────────────────────────────────────────────────────────────────────────
// API call
// ─────────────────────────────────────────────────────────────────────────────

// GenerateContent calls POST /v1/messages and returns the parsed response.
func (c *ClaudeClient) GenerateContent(ctx context.Context, req ClaudeRequest) (*ClaudeResponse, error) {
	if !c.Configured() {
		return nil, ErrClaudeUnavailable
	}
	req.Model = c.chatModel
	if req.MaxTokens <= 0 {
		req.MaxTokens = claudeDefaultMaxTokens
	}
	rawURL := claudeBaseURL + "/messages"
	var out ClaudeResponse
	if err := c.doJSON(ctx, http.MethodPost, rawURL, req, &out, "messages"); err != nil {
		return nil, err
	}
	return &out, nil
}

// ─────────────────────────────────────────────────────────────────────────────
// HTTP plumbing with one retry on 429 / 5xx
// ─────────────────────────────────────────────────────────────────────────────

func (c *ClaudeClient) doJSON(ctx context.Context, method, rawURL string, body any, out any, op string) error {
	const maxAttempts = 2
	var lastErr error
	for attempt := 1; attempt <= maxAttempts; attempt++ {
		err := c.doOnce(ctx, method, rawURL, body, out, op)
		if err == nil {
			return nil
		}
		lastErr = err
		if !isClaudeTransient(err) || attempt == maxAttempts {
			break
		}
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(time.Duration(attempt) * 2 * time.Second):
		}
	}
	if lastErr != nil {
		log.Printf("[ai][claude] op=%s giving_up_after_retries err=%v", op, lastErr)
	}
	return lastErr
}

func (c *ClaudeClient) doOnce(ctx context.Context, method, rawURL string, body any, out any, op string) error {
	var reqBody io.Reader
	if body != nil {
		raw, err := json.Marshal(body)
		if err != nil {
			return fmt.Errorf("claude: marshal request: %w", err)
		}
		reqBody = bytes.NewReader(raw)
	}
	req, err := http.NewRequestWithContext(ctx, method, rawURL, reqBody)
	if err != nil {
		return fmt.Errorf("claude: build request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("x-api-key", c.apiKey)
	req.Header.Set("anthropic-version", claudeAPIVersion)
	// prompt-caching-2024-07-31 enables cache_control on system blocks
	req.Header.Set("anthropic-beta", "prompt-caching-2024-07-31")

	resp, err := c.http.Do(req)
	if err != nil {
		log.Printf("[ai][claude] op=%s transport_error=%v", op, err)
		return fmt.Errorf("%w: %v", ErrClaudeUnavailable, err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("claude: read body: %w", err)
	}

	if resp.StatusCode >= 400 {
		snippet := truncate(string(respBody), 500)
		log.Printf("[ai][claude] op=%s status=%d body=%q", op, resp.StatusCode, snippet)
		if resp.StatusCode == http.StatusTooManyRequests || resp.StatusCode >= 500 {
			return fmt.Errorf("%w: %d: %s", ErrClaudeUnavailable, resp.StatusCode, truncate(string(respBody), 300))
		}
		return fmt.Errorf("claude: %d: %s", resp.StatusCode, truncate(string(respBody), 300))
	}
	if out == nil {
		return nil
	}
	if err := json.Unmarshal(respBody, out); err != nil {
		log.Printf("[ai][claude] op=%s json_unmarshal err=%v body_prefix=%q", op, err, truncate(string(respBody), 400))
		return fmt.Errorf("claude: decode body: %w", err)
	}
	return nil
}

func isClaudeTransient(err error) bool {
	return err != nil && errors.Is(err, ErrClaudeUnavailable)
}
