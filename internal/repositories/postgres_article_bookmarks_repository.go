package repositories

import (
	"stocet-core-service/internal/models"

	"gorm.io/gorm"
)

type PostgresArticleBookmarksRepository struct {
	db *gorm.DB
}

func NewPostgresArticleBookmarksRepository(db *gorm.DB) *PostgresArticleBookmarksRepository {
	return &PostgresArticleBookmarksRepository{db: db}
}

func (r *PostgresArticleBookmarksRepository) Add(userID, articleID string) error {
	return r.db.Exec(`
		INSERT INTO article_bookmarks (user_id, article_id)
		VALUES (?::uuid, ?::uuid)
		ON CONFLICT (user_id, article_id) DO NOTHING
	`, userID, articleID).Error
}

func (r *PostgresArticleBookmarksRepository) Remove(userID, articleID string) error {
	return r.db.Exec(`
		DELETE FROM article_bookmarks
		WHERE user_id = ?::uuid AND article_id = ?::uuid
	`, userID, articleID).Error
}

func (r *PostgresArticleBookmarksRepository) List(userID string, limit, skip int) ([]NewsArticle, int64, error) {
	if limit <= 0 {
		limit = 20
	}
	if limit > 100 {
		limit = 100
	}
	if skip < 0 {
		skip = 0
	}

	var total int64
	if err := r.db.Raw(`
		SELECT COUNT(*)
		FROM article_bookmarks
		WHERE user_id = ?::uuid
	`, userID).Scan(&total).Error; err != nil {
		return nil, 0, err
	}

	var rows []models.Article
	if err := r.db.Raw(`
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
		FROM article_bookmarks b
		JOIN articles a ON a.id = b.article_id
		WHERE b.user_id = ?::uuid
		ORDER BY b.created_at DESC
		LIMIT ? OFFSET ?
	`, userID, limit, skip).Scan(&rows).Error; err != nil {
		return nil, 0, err
	}

	articles := make([]NewsArticle, 0, len(rows))
	for _, row := range rows {
		article := mapArticleModel(row)
		article.Bookmarked = true
		articles = append(articles, article)
	}

	return articles, total, nil
}
