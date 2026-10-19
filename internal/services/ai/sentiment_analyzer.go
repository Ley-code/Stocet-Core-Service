package ai

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"time"

	"stocet-core-service/internal/repositories"
)

// SentimentResult holds the Gemini classification for one news article.
type SentimentResult struct {
	Sentiment  string   `json:"sentiment"` // "bullish" | "bearish" | "neutral"
	Tickers    []string `json:"tickers"`
	Confidence float64  `json:"confidence"`
}

// ClassifyNewsSentiment calls Gemini to classify a single article and returns
// the structured result. Returns an error only on API or parse failure — the
// caller should fall back to "neutral" in that case.
//
// validTickers is the set of real ESX ticker codes. They are given to the model and any other code it
// returns is dropped — without this the model invented plausible-looking tickers (e.g. "CBE", "AWB").
func ClassifyNewsSentiment(ctx context.Context, client *GeminiClient, title, summary string, validTickers []string) (*SentimentResult, error) {
	if !client.Configured() {
		return nil, fmt.Errorf("sentiment: gemini not configured")
	}

	articleText := title
	if s := strings.TrimSpace(summary); s != "" {
		articleText += "\n\n" + s
	}

	prompt := fmt.Sprintf(`You are a financial news analyst for the Ethiopian Securities Exchange (ESX).

Classify the sentiment of this news article and identify which ESX-listed companies it directly affects.

Article:
%s

Respond with ONLY valid JSON — no markdown, no explanation, no code fences:
{"sentiment":"neutral","tickers":[],"confidence":0.9}

Rules:
- sentiment: "bullish" (positive for investors: profit growth, expansion, dividends, strong results)
             "bearish" (negative for investors: losses, regulatory issues, scandal, poor results)
             "neutral" (general market info, policy news with unclear direction)
- tickers: ONLY codes from this list of ESX-listed securities: %s
           Include a code only if that company is directly mentioned or clearly affected. Empty array if none.
           Never output a code that is not in the list.
- confidence: 0.0 to 1.0`, articleText, tickerListForPrompt(validTickers))

	temp := 0.0
	// 150 tokens truncated the JSON mid-object on thinking models; leave generous headroom.
	maxTokens := 600
	resp, err := client.GenerateContent(ctx, GeminiRequest{
		Contents: []GeminiContent{NewUserContent(prompt)},
		GenerationConfig: &GeminiGenerationConfig{
			Temperature:      &temp,
			MaxOutputTokens:  &maxTokens,
			ResponseMimeType: "application/json",
		},
	})
	if err != nil {
		return nil, fmt.Errorf("sentiment: gemini error: %w", err)
	}

	raw := strings.TrimSpace(resp.Text())
	raw = strings.TrimPrefix(raw, "```json")
	raw = strings.TrimPrefix(raw, "```")
	raw = strings.TrimSuffix(raw, "```")
	raw = strings.TrimSpace(raw)

	var result SentimentResult
	if err := json.Unmarshal([]byte(raw), &result); err != nil {
		return nil, fmt.Errorf("sentiment: parse failed %q: %w", raw, err)
	}

	switch result.Sentiment {
	case "bullish", "bearish", "neutral":
	default:
		result.Sentiment = "neutral"
	}
	result.Tickers = filterTickers(result.Tickers, validTickers)
	if result.Confidence < 0 {
		result.Confidence = 0
	}
	if result.Confidence > 1 {
		result.Confidence = 1
	}

	return &result, nil
}

func tickerListForPrompt(valid []string) string {
	if len(valid) == 0 {
		return "(none known — return an empty array)"
	}
	return strings.Join(valid, ", ")
}

// filterTickers keeps only codes present in valid (case-insensitive), upper-cased and de-duplicated.
func filterTickers(got, valid []string) []string {
	allowed := make(map[string]bool, len(valid))
	for _, v := range valid {
		allowed[strings.ToUpper(strings.TrimSpace(v))] = true
	}
	out := []string{}
	seen := map[string]bool{}
	for _, t := range got {
		t = strings.ToUpper(strings.TrimSpace(t))
		if t != "" && allowed[t] && !seen[t] {
			seen[t] = true
			out = append(out, t)
		}
	}
	return out
}

// RunSentimentBackfillOnce fetches one batch of unanalyzed articles and
// classifies them. Returns the number of articles successfully processed.
func RunSentimentBackfillOnce(ctx context.Context, client *GeminiClient, repo repositories.SentimentBackfillRepository, batchSize int) (int, error) {
	if !client.Configured() {
		return 0, nil
	}
	if batchSize <= 0 {
		batchSize = 20
	}

	articles, err := repo.GetUnanalyzed(batchSize)
	if err != nil {
		return 0, fmt.Errorf("sentiment backfill: fetch: %w", err)
	}

	validTickers, err := repo.GetKnownTickers()
	if err != nil {
		log.Printf("[sentiment] could not load known tickers, tagging none: %v", err)
		validTickers = nil
	}

	processed := 0
	consecutiveFailures := 0
	for _, article := range articles {
		if article.ID == "" {
			continue
		}

		result, err := ClassifyNewsSentiment(ctx, client, article.Title, article.Summary, validTickers)
		if err != nil {
			// Leave the article unanalysed so a later run retries it. (Stamping "neutral" on every
			// error — including quota 429s — permanently mislabelled articles.) If the API keeps
			// failing, stop the batch rather than burning quota.
			log.Printf("[sentiment] will retry article %v later: %v", article.ID, err)
			consecutiveFailures++
			if consecutiveFailures >= 3 {
				log.Printf("[sentiment] %d failures in a row, ending this batch", consecutiveFailures)
				break
			}
			continue
		}
		consecutiveFailures = 0

		if err := repo.UpdateSentiment(ctx, article.ID, result.Sentiment, result.Tickers, &result.Confidence); err != nil {
			log.Printf("[sentiment] update failed for %v: %v", article.ID, err)
			continue
		}

		processed++

		select {
		case <-ctx.Done():
			return processed, ctx.Err()
		case <-time.After(300 * time.Millisecond):
		}
	}

	if processed > 0 {
		log.Printf("[sentiment] backfill: processed %d/%d articles", processed, len(articles))
	}

	return processed, nil
}

// StartSentimentWorker runs a periodic sentiment backfill — call once from main().
// Runs once immediately then on every interval.
func StartSentimentWorker(ctx context.Context, client *GeminiClient, repo repositories.SentimentBackfillRepository, interval time.Duration, batchSize int) {
	if !client.Configured() {
		log.Println("[sentiment] worker skipped (gemini not configured)")
		return
	}
	if interval <= 0 {
		interval = 15 * time.Minute
	}

	log.Printf("[sentiment] worker starting (interval=%s batch=%d)", interval, batchSize)
	t := time.NewTicker(interval)
	defer t.Stop()

	RunSentimentBackfillOnce(ctx, client, repo, batchSize)

	for {
		select {
		case <-ctx.Done():
			log.Println("[sentiment] worker stopping")
			return
		case <-t.C:
			RunSentimentBackfillOnce(ctx, client, repo, batchSize)
		}
	}
}
