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
	Sentiment  string   `json:"sentiment"`   // "bullish" | "bearish" | "neutral"
	Tickers    []string `json:"tickers"`
	Confidence float64  `json:"confidence"`
}

// ClassifyNewsSentiment calls Gemini to classify a single article and returns
// the structured result. Returns an error only on API or parse failure — the
// caller should fall back to "neutral" in that case.
func ClassifyNewsSentiment(ctx context.Context, client *GeminiClient, title, summary string) (*SentimentResult, error) {
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
- tickers: ESX company ticker symbols that are directly mentioned or clearly affected. Empty array if none.
           ESX tickers are short uppercase letter codes (e.g. CBE, AWB, ACSI).
- confidence: 0.0 to 1.0`, articleText)

	temp := 0.0
	maxTokens := 150
	resp, err := client.GenerateContent(ctx, GeminiRequest{
		Contents: []GeminiContent{NewUserContent(prompt)},
		GenerationConfig: &GeminiGenerationConfig{
			Temperature:     &temp,
			MaxOutputTokens: &maxTokens,
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
	if result.Tickers == nil {
		result.Tickers = []string{}
	}

	return &result, nil
}

// RunSentimentBackfillOnce fetches one batch of unanalyzed articles and
// classifies them. Returns the number of articles successfully processed.
// Safe to call concurrently — MongoDB UpdateOne is atomic per document.
func RunSentimentBackfillOnce(ctx context.Context, client *GeminiClient, repo *repositories.NewsRepository, batchSize int) (int, error) {
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

	processed := 0
	for _, article := range articles {
		if article.ID == nil {
			continue
		}

		result, err := ClassifyNewsSentiment(ctx, client, article.Title, article.Summary)
		if err != nil {
			log.Printf("[sentiment] skipping article %v: %v", article.ID, err)
			// Mark as neutral so we don't retry indefinitely on this article
			_ = repo.UpdateSentiment(ctx, *article.ID, "neutral", []string{})
			continue
		}

		if err := repo.UpdateSentiment(ctx, *article.ID, result.Sentiment, result.Tickers); err != nil {
			log.Printf("[sentiment] update failed for %v: %v", article.ID, err)
			continue
		}

		processed++

		// Pace requests to stay under Gemini rate limits
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

// StartSentimentWorker runs a periodic sentiment backfill — call from main()
// outside of Lambda. Runs once immediately then on every interval.
func StartSentimentWorker(ctx context.Context, client *GeminiClient, repo *repositories.NewsRepository, interval time.Duration, batchSize int) {
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
