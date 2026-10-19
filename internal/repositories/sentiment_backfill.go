package repositories

import "context"

// SentimentArticle is a minimal article projection used by the sentiment backfill worker.
type SentimentArticle struct {
	ID      string
	Title   string
	Summary string
}

// SentimentBackfillRepository is the minimal interface the sentiment backfill worker needs.
type SentimentBackfillRepository interface {
	GetUnanalyzed(limit int) ([]SentimentArticle, error)
	UpdateSentiment(ctx context.Context, id string, sentiment string, tickers []string, confidence *float64) error
	// GetKnownTickers returns the real ESX ticker codes, so the classifier can't tag articles with invented ones.
	GetKnownTickers() ([]string, error)
}
