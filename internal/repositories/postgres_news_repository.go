package repositories

import (
	"context"
	"fmt"
	"strings"
	"time"

	"stocet-core-service/internal/models"

	"github.com/google/uuid"
	"github.com/lib/pq"
	"gorm.io/gorm"
)

// NewsArticle is the API-facing news article shape returned by the news endpoints.
type NewsArticle struct {
	ID                  string     `json:"id,omitempty"`
	Title               string     `json:"title"`
	URL                 string     `json:"url"`
	Source              string     `json:"source"`
	PublishedDate       time.Time  `json:"published_date"`
	Content             string     `json:"content"`
	Summary             string     `json:"summary,omitempty"`
	Tags                []string   `json:"tags"`
	ImageURL            string     `json:"image_url,omitempty"`
	Sentiment           string     `json:"sentiment,omitempty"`
	SentimentTickers    []string   `json:"sentiment_tickers,omitempty"`
	SentimentAnalyzedAt *time.Time `json:"sentiment_analyzed_at,omitempty"`
	LikeCount           int64      `json:"like_count,omitempty"`
	Liked               bool       `json:"liked,omitempty"`
}

type PostgresNewsRepository struct {
	db *gorm.DB
}

func NewPostgresNewsRepository(db *gorm.DB) *PostgresNewsRepository {
	return &PostgresNewsRepository{db: db}
}

type newsSQLBuilder struct {
	conditions []string
	args       []interface{}
}

func (b *newsSQLBuilder) add(condition string, args ...interface{}) {
	b.conditions = append(b.conditions, condition)
	b.args = append(b.args, args...)
}

func (b *newsSQLBuilder) whereClause() string {
	if len(b.conditions) == 0 {
		return ""
	}
	return "WHERE " + strings.Join(b.conditions, " AND ")
}

func buildNewsSQLFilters(filters map[string]interface{}) newsSQLBuilder {
	builder := newsSQLBuilder{}

	if category, ok := filters["category"].(string); ok && category != "" {
		builder.add("? = ANY(tags)", category)
	}

	if source, ok := filters["source"].(string); ok && source != "" {
		builder.add("source = ?", source)
	}

	if startDate, ok := filters["start_date"].(time.Time); ok {
		builder.add("published_date >= ?", startDate)
	}

	if endDate, ok := filters["end_date"].(time.Time); ok {
		builder.add("published_date <= ?", endDate)
	}

	if searchText, ok := filters["search_text"].(string); ok && searchText != "" {
		builder.add(
			"to_tsvector('english', title || ' ' || coalesce(content, '')) @@ plainto_tsquery('english', ?)",
			searchText,
		)
	}

	if sentiment, ok := filters["sentiment"].(string); ok && sentiment != "" {
		if sentiment == "analyzed" {
			builder.add("(sentiment IS NOT NULL AND sentiment <> '')")
		} else {
			builder.add("sentiment = ?", sentiment)
		}
	}

	if ticker, ok := filters["ticker"].(string); ok && ticker != "" {
		builder.add("? = ANY(sentiment_tickers)", ticker)
	}

	return builder
}

func (r *PostgresNewsRepository) CountNews(filters map[string]interface{}) (int64, error) {
	builder := buildNewsSQLFilters(filters)
	query := fmt.Sprintf("SELECT COUNT(*) FROM articles %s", builder.whereClause())

	var total int64
	if err := r.db.Raw(query, builder.args...).Scan(&total).Error; err != nil {
		return 0, err
	}
	return total, nil
}

func (r *PostgresNewsRepository) GetNews(
	filters map[string]interface{},
	limit, skip int,
	cursorPublished *time.Time,
	cursorID *string,
) ([]NewsArticle, error) {
	builder := buildNewsSQLFilters(filters)

	if cursorPublished != nil && cursorID != nil {
		builder.add(
			"(published_date < ? OR (published_date = ? AND id < ?::uuid))",
			*cursorPublished,
			*cursorPublished,
			*cursorID,
		)
	}

	query := fmt.Sprintf(`
		SELECT
			id,
			title,
			url,
			source,
			published_date,
			content,
			COALESCE(summary, '') AS summary,
			COALESCE(tags, ARRAY[]::text[]) AS tags,
			COALESCE(image_url, '') AS image_url,
			COALESCE(sentiment, '') AS sentiment,
			COALESCE(sentiment_tickers, ARRAY[]::text[]) AS sentiment_tickers,
			sentiment_analyzed_at
		FROM articles
		%s
		ORDER BY published_date DESC, id DESC
	`, builder.whereClause())

	args := append([]interface{}{}, builder.args...)
	if cursorPublished == nil || cursorID == nil {
		query += " OFFSET ? LIMIT ?"
		args = append(args, skip, limit)
	} else {
		query += " LIMIT ?"
		args = append(args, limit)
	}

	var rows []models.Article
	if err := r.db.Raw(query, args...).Scan(&rows).Error; err != nil {
		return nil, err
	}

	articles := make([]NewsArticle, 0, len(rows))
	for _, row := range rows {
		articles = append(articles, mapArticleModel(row))
	}
	return articles, nil
}

func mapArticleModel(row models.Article) NewsArticle {
	tags := []string(row.Tags)
	if tags == nil {
		tags = []string{}
	}
	tickers := []string(row.SentimentTickers)
	if tickers == nil {
		tickers = []string{}
	}

	return NewsArticle{
		ID:                  row.ID.String(),
		Title:               row.Title,
		URL:                 row.URL,
		Source:              row.Source,
		PublishedDate:       row.PublishedDate,
		Content:             row.Content,
		Summary:             row.Summary,
		Tags:                tags,
		ImageURL:            row.ImageURL,
		Sentiment:           row.Sentiment,
		SentimentTickers:    tickers,
		SentimentAnalyzedAt: row.SentimentAnalyzedAt,
	}
}

func (r *PostgresNewsRepository) DeleteNews(id string) error {
	return r.db.Exec("DELETE FROM articles WHERE id = ?::uuid", id).Error
}

func (r *PostgresNewsRepository) GetUnanalyzed(limit int) ([]SentimentArticle, error) {
	if limit <= 0 {
		limit = 20
	}

	var rows []struct {
		ID      uuid.UUID
		Title   string
		Summary string
	}

	err := r.db.Raw(`
		SELECT id, title, COALESCE(summary, '') AS summary
		FROM articles
		WHERE sentiment IS NULL OR sentiment = ''
		ORDER BY published_date DESC
		LIMIT ?
	`, limit).Scan(&rows).Error
	if err != nil {
		return nil, err
	}

	articles := make([]SentimentArticle, 0, len(rows))
	for _, row := range rows {
		articles = append(articles, SentimentArticle{
			ID:      row.ID.String(),
			Title:   row.Title,
			Summary: row.Summary,
		})
	}
	return articles, nil
}

func (r *PostgresNewsRepository) UpdateSentiment(ctx context.Context, id string, sentiment string, tickers []string) error {
	if tickers == nil {
		tickers = []string{}
	}

	return r.db.WithContext(ctx).Exec(`
		UPDATE articles
		SET sentiment = ?,
		    sentiment_tickers = ?,
		    sentiment_analyzed_at = NOW()
		WHERE id = ?::uuid
	`, sentiment, pq.Array(tickers), id).Error
}
