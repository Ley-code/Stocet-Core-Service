package services

import (
	"context"
	"errors"
	"fmt"
	"log"
	"strings"

	"stocet-core-service/internal/services/ai"
)

// SystemPromptDataAssistant is the strict, hardcoded system prompt for the
// Stocet Data Assistant. Any change to this string is a behavioural change
// for every conversation, so we keep it in code (not config).
const SystemPromptDataAssistant = `You are the Stocet Data Assistant — an AI analyst for the Ethiopian Securities Exchange (ESX).

CONTEXT: The ESX is Ethiopia's securities exchange, launched in 2024. Listed companies include Ethiopian banks (CBE, Awash, Abyssinia, Dashen, etc.), insurance companies, and selected industrials. Tickers are 2–6 uppercase characters.

TOOLS AVAILABLE — always call them before answering data questions. Never guess numbers.
• SearchMarketData — semantic search over corporate news and company profiles
• GetCompanyFundamentals — structured financial snapshot (financials, EPS, board size)
• GetTechnicalIndicators — SMA, EMA, RSI, MACD, Bollinger Bands on price history
• GetMarketPrice — latest quoted price and day change for a single ticker
• GetRecentNews — most recent corporate announcements, optionally filtered by ticker
• GetMarketSummary — today's ESX market overview: top gainers and losers

RULES:
1. ALWAYS use a tool to fetch data before answering. Never cite prices or statistics from memory.
2. Do NOT provide investment advice, buy/sell recommendations, or portfolio suggestions. If asked, state you only provide data analysis.
3. When you cite a number, reference where it came from (e.g., "Per the 2023 annual report…").
4. Format for a terminal: use markdown tables for comparisons, bullet points for lists. Be concise.
5. When asked about a specific company, call both GetCompanyFundamentals and GetRecentNews before answering.
6. End with one short disclaimer only if you discussed financial figures. Never repeat disclaimers.`

// SystemPromptAdvisory is used exclusively for the Robo Advisory feature.
// It explicitly allows portfolio allocation guidance, unlike the default analyst prompt.
const SystemPromptAdvisory = `You are the Stocet Portfolio Advisor — an AI-powered portfolio guidance assistant for the Ethiopian Securities Exchange (ESX).

CONTEXT: The ESX launched in 2024. Listed companies include Ethiopian banks (CBE, Awash, Abyssinia, Dashen, etc.), insurance companies, and selected industrials. Tickers are 2–6 uppercase characters. There are approximately 20 listed securities.

TOOLS AVAILABLE — use them before forming any recommendation:
• GetMarketSummary — today's ESX market state: top gainers/losers
• GetCompanyFundamentals — financial snapshot for a company (financials, EPS, board)
• GetMarketPrice — latest price and day change for a ticker
• GetRecentNews — latest corporate announcements

YOUR TASK: Provide educational portfolio allocation guidance based on the user's risk profile. You MUST:
1. Call GetMarketSummary first to see current market conditions
2. Call GetCompanyFundamentals for the most relevant ESX companies
3. Call GetMarketPrice for companies you plan to include
4. Output a structured allocation table followed by rationale, risks, and getting-started advice

OUTPUT FORMAT — follow this exactly:
## Portfolio Allocation
| Ticker | Company | Allocation % | Rationale |
|--------|---------|-------------|-----------|

## Investment Thesis
[2–3 sentences]

## Key Risks
[3 bullet points]

## Getting Started
[1–2 sentences on executing through a Capital Market Service Provider (CMSP)]

---
*This is educational analysis based on publicly available ESX data. It does not constitute financial advice. Consult a licensed investment advisor before making investment decisions.*

RULES:
- Base every recommendation on data fetched from tools. Never guess prices or fundamentals.
- Tailor allocations to the user's stated risk profile, horizon, and goal.
- Keep total allocations to 100%.
- Include a mix of securities appropriate for the profile (conservative = stable banks; aggressive = growth-oriented).`

// ErrAIDataEngineUnavailable is the user-facing error returned when the AI
// pipeline can't service a request. Controllers map this to 503 with the
// canonical message specified in the plan.
var ErrAIDataEngineUnavailable = errors.New("the data engine is currently experiencing high load. Please try again in a moment.")

// ChatMessage is the wire format used by both the HTTP API and the service.
type ChatMessage struct {
	Role    string `json:"role"`    // "user" | "assistant"
	Content string `json:"content"`
}

// AIService orchestrates Gemini chat with tool-calling: RAG search and
// technical-analysis lookups. It owns nothing concurrent; one Chat call =
// one synchronous Gemini interaction (with up to 3 tool rounds).
type AIService struct {
	client       *ai.GeminiClient
	rag          *ai.RAGService
	ta           *ai.TechnicalAnalysisService
	maxTurns     int
	maxChars     int
	maxToolRounds int
}

// NewAIService wires the chat client + tools. maxTurns and maxChars cap
// how much history we ship per request; pass 0 to use sane defaults.
func NewAIService(client *ai.GeminiClient, rag *ai.RAGService, ta *ai.TechnicalAnalysisService, maxTurns, maxChars int) *AIService {
	if maxTurns <= 0 {
		maxTurns = 12
	}
	if maxChars <= 0 {
		maxChars = 16000
	}
	return &AIService{
		client:        client,
		rag:           rag,
		ta:            ta,
		maxTurns:      maxTurns,
		maxChars:      maxChars,
		maxToolRounds: 3,
	}
}

// Configured reports whether the underlying Gemini client has an API key.
// Controllers can fast-fail with the graceful error when this is false.
func (s *AIService) Configured() bool {
	return s != nil && s.client != nil && s.client.Configured()
}

// Chat runs one conversation turn. `history` includes the new user message
// at the end. Returns the assistant's final text or ErrAIDataEngineUnavailable.
func systemPromptForMode(mode string) string {
	if mode == "advisory" {
		return SystemPromptAdvisory
	}
	return SystemPromptDataAssistant
}

func (s *AIService) Chat(ctx context.Context, history []ChatMessage, mode string) (string, error) {
	if !s.Configured() {
		return "", ErrAIDataEngineUnavailable
	}

	trimmed := truncateHistory(history, s.maxTurns, s.maxChars)
	if len(trimmed) == 0 {
		log.Printf("[ai][chat] abort reason=empty_history_after_truncation max_turns=%d max_chars=%d", s.maxTurns, s.maxChars)
		return "", errors.New("ai: empty history")
	}

	chars := 0
	for _, m := range trimmed {
		chars += len(m.Content)
	}
	log.Printf("[ai][chat] gemini_request model=%s turns=%d content_chars=%d max_turns=%d max_chars=%d",
		s.client.ChatModel(), len(trimmed), chars, s.maxTurns, s.maxChars)

	contents := make([]ai.GeminiContent, 0, len(trimmed))
	for _, m := range trimmed {
		switch strings.ToLower(strings.TrimSpace(m.Role)) {
		case "user":
			contents = append(contents, ai.NewUserContent(m.Content))
		case "assistant", "model":
			contents = append(contents, ai.NewModelContent(m.Content))
		default:
			// Unknown roles are treated as user — defensive default.
			contents = append(contents, ai.NewUserContent(m.Content))
		}
	}

	tools := s.buildTools()
	req := ai.GeminiRequest{
		SystemInstruction: ai.NewSystemInstruction(systemPromptForMode(mode)),
		Contents:          contents,
		Tools:             tools,
	}

	for round := 0; round <= s.maxToolRounds; round++ {
		resp, err := s.client.GenerateContent(ctx, req)
		if err != nil {
			log.Printf("[ai][chat] generateContent_failed round=%d err=%v", round, err)
			return "", ErrAIDataEngineUnavailable
		}

		calls := resp.FunctionCalls()
		if len(calls) == 0 {
			text := strings.TrimSpace(resp.Text())
			if text == "" {
				log.Printf("[ai][chat] empty_model_text round=%d %s", round, summarizeGeminiResponse(resp))
				return "", ErrAIDataEngineUnavailable
			}
			log.Printf("[ai][chat] complete round=%d reply_chars=%d", round, len(text))
			return text, nil
		}
		if round == s.maxToolRounds {
			// We've hit the tool-call cap; return whatever text the model
			// produced alongside the calls, or the graceful error.
			text := strings.TrimSpace(resp.Text())
			if text != "" {
				log.Printf("[ai][chat] complete_after_tool_cap round=%d reply_chars=%d", round, len(text))
				return text, nil
			}
			log.Printf("[ai][chat] empty_model_text_after_tool_cap round=%d %s", round, summarizeGeminiResponse(resp))
			return "", ErrAIDataEngineUnavailable
		}

		toolNames := make([]string, len(calls))
		for i := range calls {
			toolNames[i] = calls[i].Name
		}
		log.Printf("[ai][chat] tool_round=%d tools=%v", round, toolNames)

		// Echo the model's tool-call turn back into history, then append our
		// function responses, then re-request.
		fc := resp.FirstCandidate()
		if fc == nil {
			log.Printf("[ai][chat] nil_first_candidate round=%d candidates=%d", round, len(resp.Candidates))
			return "", ErrAIDataEngineUnavailable
		}
		req.Contents = append(req.Contents, fc.Content)
		for _, call := range calls {
			result := s.dispatchToolCall(ctx, call)
			if errVal, ok := result["error"].(string); ok && errVal != "" {
				log.Printf("[ai][tool] name=%s tool_error=%q", call.Name, errVal)
			} else {
				log.Printf("[ai][tool] name=%s ok", call.Name)
			}
			req.Contents = append(req.Contents, ai.NewFunctionResponseContent(call.Name, result))
		}
	}

	log.Printf("[ai][chat] exhausted_tool_loop max_rounds=%d", s.maxToolRounds)
	return "", ErrAIDataEngineUnavailable
}

// StreamEventKind enumerates the SSE event types emitted by ChatStream.
type StreamEventKind string

const (
	StreamEventTool  StreamEventKind = "tool"
	StreamEventText  StreamEventKind = "text"
	StreamEventError StreamEventKind = "error"
)

// StreamEvent is one SSE payload emitted by ChatStream.
type StreamEvent struct {
	Event  StreamEventKind `json:"event"`
	Tool   string          `json:"tool,omitempty"`   // tool events
	Status string          `json:"status,omitempty"` // "start" | "done"
	Text   string          `json:"text,omitempty"`   // text event
	Error  string          `json:"error,omitempty"`  // error event
}

// StreamEmitter is the callback ChatStream uses to push events to the HTTP layer.
type StreamEmitter func(StreamEvent)

// ChatStream behaves like Chat but emits progress events (tool start/done) and
// the final text via the supplied emitter. This enables real-time SSE streaming
// on the HTTP layer without needing the Gemini streaming API.
func (s *AIService) ChatStream(ctx context.Context, history []ChatMessage, mode string, emit StreamEmitter) error {
	if !s.Configured() {
		emit(StreamEvent{Event: StreamEventError, Error: ErrAIDataEngineUnavailable.Error()})
		return ErrAIDataEngineUnavailable
	}

	trimmed := truncateHistory(history, s.maxTurns, s.maxChars)
	if len(trimmed) == 0 {
		emit(StreamEvent{Event: StreamEventError, Error: "empty history"})
		return errors.New("ai: empty history")
	}

	contents := make([]ai.GeminiContent, 0, len(trimmed))
	for _, m := range trimmed {
		switch strings.ToLower(strings.TrimSpace(m.Role)) {
		case "user":
			contents = append(contents, ai.NewUserContent(m.Content))
		case "assistant", "model":
			contents = append(contents, ai.NewModelContent(m.Content))
		default:
			contents = append(contents, ai.NewUserContent(m.Content))
		}
	}

	tools := s.buildTools()
	req := ai.GeminiRequest{
		SystemInstruction: ai.NewSystemInstruction(systemPromptForMode(mode)),
		Contents:          contents,
		Tools:             tools,
	}

	for round := 0; round <= s.maxToolRounds; round++ {
		resp, err := s.client.GenerateContent(ctx, req)
		if err != nil {
			log.Printf("[ai][stream] generateContent_failed round=%d err=%v", round, err)
			emit(StreamEvent{Event: StreamEventError, Error: ErrAIDataEngineUnavailable.Error()})
			return ErrAIDataEngineUnavailable
		}

		calls := resp.FunctionCalls()
		if len(calls) == 0 {
			text := strings.TrimSpace(resp.Text())
			if text == "" {
				log.Printf("[ai][stream] empty_model_text round=%d", round)
				emit(StreamEvent{Event: StreamEventError, Error: ErrAIDataEngineUnavailable.Error()})
				return ErrAIDataEngineUnavailable
			}
			log.Printf("[ai][stream] complete round=%d text_chars=%d", round, len(text))
			emit(StreamEvent{Event: StreamEventText, Text: text})
			return nil
		}

		if round == s.maxToolRounds {
			text := strings.TrimSpace(resp.Text())
			if text != "" {
				emit(StreamEvent{Event: StreamEventText, Text: text})
				return nil
			}
			emit(StreamEvent{Event: StreamEventError, Error: ErrAIDataEngineUnavailable.Error()})
			return ErrAIDataEngineUnavailable
		}

		fc := resp.FirstCandidate()
		if fc == nil {
			emit(StreamEvent{Event: StreamEventError, Error: ErrAIDataEngineUnavailable.Error()})
			return ErrAIDataEngineUnavailable
		}
		req.Contents = append(req.Contents, fc.Content)

		for _, call := range calls {
			log.Printf("[ai][stream][tool] round=%d name=%s status=start", round, call.Name)
			emit(StreamEvent{Event: StreamEventTool, Tool: call.Name, Status: "start"})
			result := s.dispatchToolCall(ctx, call)
			if errVal, ok := result["error"].(string); ok && errVal != "" {
				log.Printf("[ai][stream][tool] round=%d name=%s tool_error=%q", round, call.Name, errVal)
			}
			emit(StreamEvent{Event: StreamEventTool, Tool: call.Name, Status: "done"})
			req.Contents = append(req.Contents, ai.NewFunctionResponseContent(call.Name, result))
		}
	}

	emit(StreamEvent{Event: StreamEventError, Error: ErrAIDataEngineUnavailable.Error()})
	return ErrAIDataEngineUnavailable
}

// summarizeGeminiResponse builds a short, log-safe summary when the model
// returns no text (blocked, malformed, or empty candidates).
func summarizeGeminiResponse(resp *ai.GeminiResponse) string {
	if resp == nil {
		return "response=nil"
	}
	n := len(resp.Candidates)
	if n == 0 {
		return "candidates=0"
	}
	c := resp.Candidates[0]
	parts := len(c.Content.Parts)
	return fmt.Sprintf("candidates=%d finish_reason=%q parts=%d", n, c.FinishReason, parts)
}

// dispatchToolCall routes a single tool call to the right service.
// Errors are surfaced to the LLM as `{ "error": "..." }` so it can apologize
// gracefully instead of failing the whole chat.
func (s *AIService) dispatchToolCall(ctx context.Context, call ai.GeminiFunctionCall) map[string]interface{} {
	switch call.Name {
	case "SearchMarketData":
		query, _ := call.Args["query"].(string)
		hits, err := s.rag.SearchMarketData(ctx, query, 6)
		if err != nil {
			return map[string]interface{}{"error": err.Error()}
		}
		return map[string]interface{}{"results": hits, "count": len(hits)}

	case "GetCompanyFundamentals":
		ticker, _ := call.Args["ticker"].(string)
		snap, err := s.rag.GetCompanyFundamentals(ctx, ticker)
		if err != nil {
			return map[string]interface{}{"error": err.Error()}
		}
		return map[string]interface{}{"fundamentals": snap}

	case "GetTechnicalIndicators":
		ticker, _ := call.Args["ticker"].(string)
		indicator, _ := call.Args["indicator"].(string)
		// JSON numbers decode to float64 in Go.
		periodF, _ := call.Args["period"].(float64)
		period := int(periodF)
		res, err := s.ta.ComputeIndicator(ctx, ticker, indicator, period)
		if err != nil {
			return map[string]interface{}{"error": err.Error()}
		}
		return map[string]interface{}{"indicator": res}

	case "GetMarketPrice":
		ticker, _ := call.Args["ticker"].(string)
		result, err := s.rag.GetLatestPrice(ctx, ticker)
		if err != nil {
			return map[string]interface{}{"error": err.Error()}
		}
		return map[string]interface{}{"price": result}

	case "GetRecentNews":
		ticker, _ := call.Args["ticker"].(string)
		limitF, _ := call.Args["limit"].(float64)
		limit := int(limitF)
		results, err := s.rag.GetRecentNews(ctx, ticker, limit)
		if err != nil {
			return map[string]interface{}{"error": err.Error()}
		}
		return map[string]interface{}{"news": results, "count": len(results)}

	case "GetMarketSummary":
		summary, err := s.rag.GetMarketSummary(ctx)
		if err != nil {
			return map[string]interface{}{"error": err.Error()}
		}
		return map[string]interface{}{"summary": summary}
	}

	return map[string]interface{}{"error": "unknown tool: " + call.Name}
}

// buildTools declares the function signatures the model is allowed to call.
func (s *AIService) buildTools() []ai.GeminiTool {
	return []ai.GeminiTool{{
		FunctionDeclarations: []ai.GeminiFunctionDeclaration{
			{
				Name:        "SearchMarketData",
				Description: "Semantic search over Ethiopian capital-markets corporate news and company profiles. Use this to retrieve relevant context before answering questions about specific companies, news, or sectors.",
				Parameters: ai.GeminiSchema{
					Type: "object",
					Properties: map[string]ai.GeminiSchema{
						"query": {Type: "string", Description: "A natural-language search query (e.g., 'CBE Q3 dividend announcement')."},
					},
					Required: []string{"query"},
				},
			},
			{
				Name:        "GetCompanyFundamentals",
				Description: "Retrieve a structured snapshot of a single listed company by ticker: profile, latest financial reports, latest EPS, and board size.",
				Parameters: ai.GeminiSchema{
					Type: "object",
					Properties: map[string]ai.GeminiSchema{
						"ticker": {Type: "string", Description: "Ticker symbol on the Ethiopian Securities Exchange (e.g., 'CBE', 'AWBE')."},
					},
					Required: []string{"ticker"},
				},
			},
			{
				Name:        "GetTechnicalIndicators",
				Description: "Compute a technical indicator for a ticker over its historical close prices. Use this whenever the user asks about technical signals such as moving averages, RSI, MACD, or Bollinger Bands.",
				Parameters: ai.GeminiSchema{
					Type: "object",
					Properties: map[string]ai.GeminiSchema{
						"ticker": {Type: "string", Description: "Ticker symbol."},
						"indicator": {
							Type:        "string",
							Description: "Which indicator to compute.",
							Enum:        ai.SupportedIndicators,
						},
						"period": {
							Type:        "integer",
							Description: "Window length. Defaults: SMA/Bollinger=20, EMA=12, RSI=14. Ignored for MACD (uses 12/26/9).",
						},
					},
					Required: []string{"ticker", "indicator"},
				},
			},
			{
				Name:        "GetMarketPrice",
				Description: "Get the latest quoted price, daily percent change, and variation for a specific ticker on the ESX. Use this whenever the user asks about the current or recent price of a company.",
				Parameters: ai.GeminiSchema{
					Type: "object",
					Properties: map[string]ai.GeminiSchema{
						"ticker": {Type: "string", Description: "Ticker symbol (e.g., 'CBE', 'AWBE')."},
					},
					Required: []string{"ticker"},
				},
			},
			{
				Name:        "GetRecentNews",
				Description: "Fetch the most recent corporate news and announcements from the ESX. Optionally filter by ticker. Use this when the user asks about recent events, announcements, or news for a company or the market broadly.",
				Parameters: ai.GeminiSchema{
					Type: "object",
					Properties: map[string]ai.GeminiSchema{
						"ticker": {Type: "string", Description: "Ticker symbol to filter news. Leave empty for all recent ESX news."},
						"limit":  {Type: "integer", Description: "Number of news items to return (1–15, default 8)."},
					},
				},
			},
			{
				Name:        "GetMarketSummary",
				Description: "Get today's ESX market overview: top 5 gainers and top 5 losers by percent change, plus total number of listed tickers. Use this when the user asks how the market is doing, what moved today, or wants a market overview.",
				Parameters: ai.GeminiSchema{
					Type:       "object",
					Properties: map[string]ai.GeminiSchema{},
				},
			},
		},
	}}
}

// truncateHistory keeps the most recent maxTurns messages, then drops oldest
// messages (other than the very last one — the current user query) until the
// total content size is at or below maxChars.
func truncateHistory(history []ChatMessage, maxTurns, maxChars int) []ChatMessage {
	if len(history) == 0 {
		return nil
	}
	if maxTurns > 0 && len(history) > maxTurns {
		history = history[len(history)-maxTurns:]
	}
	if maxChars <= 0 {
		return append([]ChatMessage(nil), history...)
	}

	out := append([]ChatMessage(nil), history...)
	total := 0
	for _, m := range out {
		total += len(m.Content)
	}
	for total > maxChars && len(out) > 1 {
		total -= len(out[0].Content)
		out = out[1:]
	}
	// If the single-remaining last message is itself larger than maxChars, hard-truncate it.
	if len(out) == 1 && len(out[0].Content) > maxChars {
		out[0].Content = out[0].Content[:maxChars]
	}
	return out
}
