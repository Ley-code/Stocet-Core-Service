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
	SentimentConfidence *float64   `json:"sentiment_confidence,omitempty"`
	SentimentTickers    []string   `json:"sentiment_tickers,omitempty"`
	SentimentAnalyzedAt *time.Time `json:"sentiment_analyzed_at,omitempty"`
	LikeCount           int64      `json:"like_count"`
	CommentCount        int64      `json:"comment_count"`
	ReadingTimeMin      int        `json:"reading_time_min"`
	Liked               bool       `json:"liked,omitempty"`
	Bookmarked          bool       `json:"bookmarked,omitempty"`
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

	// category / source / sentiment accept comma-separated lists (OR within a field, AND across fields).
	if category, ok := filters["category"].(string); ok && category != "" {
		builder.add("tags && ?::text[]", pq.Array(splitList(category)))
	}

	if source, ok := filters["source"].(string); ok && source != "" {
		builder.add("source = ANY(?::text[])", pq.Array(splitList(source)))
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
		switch sentiment {
		case "analyzed":
			builder.add("(sentiment IS NOT NULL AND sentiment <> '')")
		case "unanalyzed":
			builder.add("(sentiment IS NULL OR sentiment = '')")
		default:
			builder.add("sentiment = ANY(?::text[])", pq.Array(splitList(sentiment)))
		}
	}

	if ticker, ok := filters["ticker"].(string); ok && ticker != "" {
		builder.add("? = ANY(sentiment_tickers)", ticker)
	}

	return builder
}

func splitList(v string) []string {
	parts := strings.Split(v, ",")
	out := make([]string, 0, len(parts))
	for _, p := range parts {
		if p = strings.TrimSpace(p); p != "" {
			out = append(out, p)
		}
	}
	return out
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

	if sort, _ := filters["sort"].(string); sort == "most_liked" || sort == "most_discussed" {
		return r.getNewsByEngagement(builder, sort, limit, skip)
	}

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
			sentiment_confidence,
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

// getNewsByEngagement orders by like or comment volume (newest first on ties). Offset pagination only.
func (r *PostgresNewsRepository) getNewsByEngagement(builder newsSQLBuilder, sort string, limit, skip int) ([]NewsArticle, error) {
	order := "likes DESC, comments DESC"
	if sort == "most_discussed" {
		order = "comments DESC, likes DESC"
	}
	query := fmt.Sprintf(`
		SELECT
			a.id, a.title, a.url, a.source, a.published_date, a.content,
			COALESCE(a.summary, '') AS summary,
			COALESCE(a.tags, ARRAY[]::text[]) AS tags,
			COALESCE(a.image_url, '') AS image_url,
			COALESCE(a.sentiment, '') AS sentiment,
			a.sentiment_confidence,
			COALESCE(a.sentiment_tickers, ARRAY[]::text[]) AS sentiment_tickers,
			a.sentiment_analyzed_at
		FROM (SELECT * FROM articles %s) a
		LEFT JOIN (SELECT article_id, COUNT(*) AS n FROM article_likes GROUP BY article_id) l ON l.article_id = a.id
		LEFT JOIN (SELECT article_id, COUNT(*) AS n FROM article_comments GROUP BY article_id) c ON c.article_id = a.id
		CROSS JOIN LATERAL (SELECT COALESCE(l.n, 0) AS likes, COALESCE(c.n, 0) AS comments) e
		ORDER BY %s, a.published_date DESC, a.id DESC
		OFFSET ? LIMIT ?
	`, builder.whereClause(), order)

	args := append(append([]interface{}{}, builder.args...), skip, limit)
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

// NewsFacets summarises the recent news mix for filter UIs and the metrics strip.
type NewsFacets struct {
	Total      int64            `json:"total"`
	Days       int              `json:"days"`
	Sources    []FacetCount     `json:"sources"`
	Categories []FacetCount     `json:"categories"`
	Tickers    []FacetCount     `json:"tickers"`
	Sentiment  map[string]int64 `json:"sentiment"`
	Daily      []DailyNewsCount `json:"daily"`
}

type FacetCount struct {
	Name  string `json:"name"`
	Count int64  `json:"count"`
}

type DailyNewsCount struct {
	Date    string `json:"date"`
	Total   int64  `json:"total"`
	Bullish int64  `json:"bullish"`
	Neutral int64  `json:"neutral"`
	Bearish int64  `json:"bearish"`
}

func (r *PostgresNewsRepository) GetNewsFacets(days int) (*NewsFacets, error) {
	if days <= 0 || days > 365 {
		days = 30
	}
	window := "published_date >= NOW() - make_interval(days => ?)"
	f := &NewsFacets{Days: days, Sentiment: map[string]int64{"bullish": 0, "neutral": 0, "bearish": 0, "unanalyzed": 0}}

	if err := r.db.Raw("SELECT COUNT(*) FROM articles WHERE "+window, days).Scan(&f.Total).Error; err != nil {
		return nil, err
	}
	if err := r.db.Raw(`SELECT source AS name, COUNT(*) AS count FROM articles WHERE `+window+` GROUP BY source ORDER BY count DESC`, days).Scan(&f.Sources).Error; err != nil {
		return nil, err
	}
	if err := r.db.Raw(`SELECT t AS name, COUNT(*) AS count FROM articles, unnest(tags) t WHERE `+window+` GROUP BY t ORDER BY count DESC`, days).Scan(&f.Categories).Error; err != nil {
		return nil, err
	}
	if err := r.db.Raw(`SELECT t AS name, COUNT(*) AS count FROM articles, unnest(sentiment_tickers) t WHERE `+window+` GROUP BY t ORDER BY count DESC LIMIT 12`, days).Scan(&f.Tickers).Error; err != nil {
		return nil, err
	}

	var sent []FacetCount
	if err := r.db.Raw(`SELECT COALESCE(NULLIF(sentiment, ''), 'unanalyzed') AS name, COUNT(*) AS count FROM articles WHERE `+window+` GROUP BY 1`, days).Scan(&sent).Error; err != nil {
		return nil, err
	}
	for _, s := range sent {
		f.Sentiment[s.Name] = s.Count
	}

	if err := r.db.Raw(`
		SELECT to_char(d::date, 'YYYY-MM-DD') AS date,
		       COUNT(a.id) AS total,
		       COUNT(a.id) FILTER (WHERE a.sentiment = 'bullish') AS bullish,
		       COUNT(a.id) FILTER (WHERE a.sentiment = 'neutral') AS neutral,
		       COUNT(a.id) FILTER (WHERE a.sentiment = 'bearish') AS bearish
		FROM generate_series((NOW() - make_interval(days => ?))::date, NOW()::date, interval '1 day') d
		LEFT JOIN articles a ON a.published_date::date = d::date
		GROUP BY d ORDER BY d`, days-1).Scan(&f.Daily).Error; err != nil {
		return nil, err
	}
	return f, nil
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
		SentimentConfidence: row.SentimentConfidence,
		SentimentTickers:    tickers,
		SentimentAnalyzedAt: row.SentimentAnalyzedAt,
		ReadingTimeMin:      readingTimeMinutes(row.Content),
	}
}

// readingTimeMinutes estimates reading time at ~200 words per minute (minimum 1).
func readingTimeMinutes(content string) int {
	words := len(strings.Fields(content))
	if words == 0 {
		return 0
	}
	if m := (words + 199) / 200; m > 1 {
		return m
	}
	return 1
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

// GetKnownTickers returns ticker codes that have a recent price on record plus any with a company profile.
func (r *PostgresNewsRepository) GetKnownTickers() ([]string, error) {
	var tickers []string
	err := r.db.Raw(`
		SELECT DISTINCT upper(ticker_symbol) FROM (
			SELECT ticker_symbol FROM market_prices WHERE timestamp > NOW() - interval '120 days'
			UNION
			SELECT ticker_symbol FROM company_profiles
		) t
		WHERE ticker_symbol IS NOT NULL AND ticker_symbol <> ''
		ORDER BY 1`).Scan(&tickers).Error
	return tickers, err
}

func (r *PostgresNewsRepository) UpdateSentiment(
	ctx context.Context,
	id string,
	sentiment string,
	tickers []string,
	confidence *float64,
) error {
	if tickers == nil {
		tickers = []string{}
	}

	return r.db.WithContext(ctx).Exec(`
		UPDATE articles
		SET sentiment = ?,
		    sentiment_confidence = ?,
		    sentiment_tickers = ?,
		    sentiment_analyzed_at = NOW()
		WHERE id = ?::uuid
	`, sentiment, confidence, pq.Array(tickers), id).Error
}

func (r *PostgresNewsRepository) GetRelatedArticles(articleID string, limit int) ([]NewsArticle, error) {
	if limit <= 0 {
		limit = 6
	}
	if limit > 25 {
		limit = 25
	}

	query := `
		WITH target AS (
			SELECT
				COALESCE(tags, ARRAY[]::text[]) AS tags,
				COALESCE(sentiment_tickers, ARRAY[]::text[]) AS sentiment_tickers
			FROM articles
			WHERE id = ?::uuid
		)
		SELECT
			a.id,
			a.title,
			a.url,
			a.source,
			a.published_date,
			a.content,
			COALESCE(a.summary, '') AS summary,
			COALESCE(a.tags, ARRAY[]::text[]) AS tags,
			COALESCE(a.image_url, '') AS image_url,
			COALESCE(a.sentiment, '') AS sentiment,
			a.sentiment_confidence,
			COALESCE(a.sentiment_tickers, ARRAY[]::text[]) AS sentiment_tickers,
			a.sentiment_analyzed_at
		FROM articles a
		CROSS JOIN target t
		WHERE a.id <> ?::uuid
		  AND (
		  	a.tags && t.tags
		  	OR a.sentiment_tickers && t.sentiment_tickers
		  )
		ORDER BY
			(
				COALESCE((
					SELECT COUNT(*)::int
					FROM unnest(COALESCE(a.sentiment_tickers, ARRAY[]::text[])) x
					WHERE x = ANY(t.sentiment_tickers)
				), 0) * 2
				+
				COALESCE((
					SELECT COUNT(*)::int
					FROM unnest(COALESCE(a.tags, ARRAY[]::text[])) y
					WHERE y = ANY(t.tags)
				), 0)
			) DESC,
			a.published_date DESC
		LIMIT ?
	`

	var rows []models.Article
	if err := r.db.Raw(query, articleID, articleID, limit).Scan(&rows).Error; err != nil {
		return nil, err
	}

	articles := make([]NewsArticle, 0, len(rows))
	for _, row := range rows {
		articles = append(articles, mapArticleModel(row))
	}
	return articles, nil
}

func (r *PostgresNewsRepository) GetTrendingArticles(limit, hours int) ([]NewsArticle, error) {
	if limit <= 0 {
		limit = 10
	}
	if limit > 50 {
		limit = 50
	}
	if hours <= 0 {
		hours = 24
	}

	query := `
		SELECT
			a.id,
			a.title,
			a.url,
			a.source,
			a.published_date,
			a.content,
			COALESCE(a.summary, '') AS summary,
			COALESCE(a.tags, ARRAY[]::text[]) AS tags,
			COALESCE(a.image_url, '') AS image_url,
			COALESCE(a.sentiment, '') AS sentiment,
			a.sentiment_confidence,
			COALESCE(a.sentiment_tickers, ARRAY[]::text[]) AS sentiment_tickers,
			a.sentiment_analyzed_at,
			COALESCE(COUNT(al.user_id), 0) AS like_count
		FROM articles a
		LEFT JOIN article_likes al
			ON al.article_id = a.id
		   AND al.created_at >= NOW() - make_interval(hours => ?)
		WHERE a.published_date >= NOW() - make_interval(hours => ?)
		GROUP BY a.id
		ORDER BY like_count DESC, a.published_date DESC
		LIMIT ?
	`

	var rows []struct {
		models.Article
		LikeCount int64 `gorm:"column:like_count"`
	}
	if err := r.db.Raw(query, hours, hours, limit).Scan(&rows).Error; err != nil {
		return nil, err
	}

	articles := make([]NewsArticle, 0, len(rows))
	for _, row := range rows {
		article := mapArticleModel(row.Article)
		article.LikeCount = row.LikeCount
		articles = append(articles, article)
	}
	return articles, nil
}

func (r *PostgresNewsRepository) GetForYouArticles(userID string, tickers []string, limit int) ([]NewsArticle, error) {
	if limit <= 0 {
		limit = 20
	}
	if limit > 100 {
		limit = 100
	}
	if len(tickers) == 0 {
		return []NewsArticle{}, nil
	}

	query := `
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
			sentiment_confidence,
			COALESCE(sentiment_tickers, ARRAY[]::text[]) AS sentiment_tickers,
			sentiment_analyzed_at
		FROM articles
		WHERE sentiment_tickers && ?::text[]
		ORDER BY published_date DESC
		LIMIT ?
	`

	_ = userID

	var rows []models.Article
	if err := r.db.Raw(query, pq.Array(tickers), limit).Scan(&rows).Error; err != nil {
		return nil, err
	}

	articles := make([]NewsArticle, 0, len(rows))
	for _, row := range rows {
		articles = append(articles, mapArticleModel(row))
	}
	return articles, nil
}
