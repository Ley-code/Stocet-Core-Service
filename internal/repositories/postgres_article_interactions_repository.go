package repositories

import (
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"
	"gorm.io/gorm"
)

type ArticleLikeSummary struct {
	Count int64
	Liked bool
}

type ArticleComment struct {
	ID        string    `json:"id"`
	ArticleID string    `json:"article_id"`
	UserID    string    `json:"user_id"`
	UserName  string    `json:"user_name"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
}

type PostgresArticleInteractionsRepository struct {
	db *gorm.DB
}

func NewPostgresArticleInteractionsRepository(db *gorm.DB) *PostgresArticleInteractionsRepository {
	return &PostgresArticleInteractionsRepository{db: db}
}

func (r *PostgresArticleInteractionsRepository) ToggleLike(articleID, userID string) (bool, error) {
	var exists bool
	err := r.db.Raw(`
		SELECT EXISTS(
			SELECT 1 FROM article_likes
			WHERE article_id = ?::uuid AND user_id = ?::uuid
		)
	`, articleID, userID).Scan(&exists).Error
	if err != nil {
		return false, err
	}

	if exists {
		err = r.db.Exec(`
			DELETE FROM article_likes
			WHERE article_id = ?::uuid AND user_id = ?::uuid
		`, articleID, userID).Error
		return false, err
	}

	err = r.db.Exec(`
		INSERT INTO article_likes (article_id, user_id)
		VALUES (?::uuid, ?::uuid)
		ON CONFLICT (article_id, user_id) DO NOTHING
	`, articleID, userID).Error
	return true, err
}

func (r *PostgresArticleInteractionsRepository) GetLikeStatus(articleID, userID string) (ArticleLikeSummary, error) {
	var summary ArticleLikeSummary
	query := `
		SELECT COUNT(*) AS count,
		       COALESCE(bool_or(user_id = ?::uuid), false) AS liked
		FROM article_likes
		WHERE article_id = ?::uuid
	`
	args := []interface{}{userID, articleID}
	if userID == "" {
		query = `
			SELECT COUNT(*) AS count, false AS liked
			FROM article_likes
			WHERE article_id = ?::uuid
		`
		args = []interface{}{articleID}
	}

	err := r.db.Raw(query, args...).Scan(&summary).Error
	return summary, err
}

type ArticleLikeCountRow struct {
	ArticleID string
	Count     int64
	Liked     bool
}

func (r *PostgresArticleInteractionsRepository) GetLikeCounts(articleIDs []string, userID string) (map[string]ArticleLikeCountRow, error) {
	result := make(map[string]ArticleLikeCountRow, len(articleIDs))
	if len(articleIDs) == 0 {
		return result, nil
	}

	var rows []ArticleLikeCountRow
	query := `
		SELECT article_id::text AS article_id,
		       COUNT(*) AS count,
		       COALESCE(bool_or(user_id = ?::uuid), false) AS liked
		FROM article_likes
		WHERE article_id = ANY(?::uuid[])
		GROUP BY article_id
	`
	args := []interface{}{userID, pq.Array(articleIDs)}
	if userID == "" {
		query = `
			SELECT article_id::text AS article_id,
			       COUNT(*) AS count,
			       false AS liked
			FROM article_likes
			WHERE article_id = ANY(?::uuid[])
			GROUP BY article_id
		`
		args = []interface{}{pq.Array(articleIDs)}
	}

	if err := r.db.Raw(query, args...).Scan(&rows).Error; err != nil {
		return nil, err
	}

	for _, row := range rows {
		result[row.ArticleID] = row
	}
	return result, nil
}

func (r *PostgresArticleInteractionsRepository) CreateComment(articleID, userID, userName, content string) (*ArticleComment, error) {
	var row struct {
		ID        uuid.UUID
		ArticleID uuid.UUID
		UserID    uuid.UUID
		UserName  string
		Content   string
		CreatedAt time.Time
	}

	err := r.db.Raw(`
		INSERT INTO article_comments (article_id, user_id, user_name, content)
		VALUES (?::uuid, ?::uuid, ?, ?)
		RETURNING id, article_id, user_id, user_name, content, created_at
	`, articleID, userID, userName, content).Scan(&row).Error
	if err != nil {
		return nil, err
	}

	return &ArticleComment{
		ID:        row.ID.String(),
		ArticleID: row.ArticleID.String(),
		UserID:    row.UserID.String(),
		Content:   row.Content,
		UserName:  row.UserName,
		CreatedAt: row.CreatedAt,
	}, nil
}

func (r *PostgresArticleInteractionsRepository) GetComments(articleID string, limit, skip int) ([]ArticleComment, int64, error) {
	var total int64
	if err := r.db.Raw(`
		SELECT COUNT(*) FROM article_comments WHERE article_id = ?::uuid
	`, articleID).Scan(&total).Error; err != nil {
		return nil, 0, err
	}

	var rows []struct {
		ID        uuid.UUID
		ArticleID uuid.UUID
		UserID    uuid.UUID
		UserName  string
		Content   string
		CreatedAt time.Time
	}

	err := r.db.Raw(`
		SELECT id, article_id, user_id, user_name, content, created_at
		FROM article_comments
		WHERE article_id = ?::uuid
		ORDER BY created_at DESC
		LIMIT ? OFFSET ?
	`, articleID, limit, skip).Scan(&rows).Error
	if err != nil {
		return nil, 0, err
	}

	comments := make([]ArticleComment, 0, len(rows))
	for _, row := range rows {
		comments = append(comments, ArticleComment{
			ID:        row.ID.String(),
			ArticleID: row.ArticleID.String(),
			UserID:    row.UserID.String(),
			UserName:  row.UserName,
			Content:   row.Content,
			CreatedAt: row.CreatedAt,
		})
	}
	return comments, total, nil
}

func (r *PostgresArticleInteractionsRepository) GetCommentByID(commentID string) (*ArticleComment, error) {
	var row struct {
		ID        uuid.UUID
		ArticleID uuid.UUID
		UserID    uuid.UUID
		UserName  string
		Content   string
		CreatedAt time.Time
	}

	err := r.db.Raw(`
		SELECT id, article_id, user_id, user_name, content, created_at
		FROM article_comments
		WHERE id = ?::uuid
	`, commentID).Scan(&row).Error
	if err != nil {
		return nil, err
	}
	if row.ID == uuid.Nil {
		return nil, gorm.ErrRecordNotFound
	}

	return &ArticleComment{
		ID:        row.ID.String(),
		ArticleID: row.ArticleID.String(),
		UserID:    row.UserID.String(),
		UserName:  row.UserName,
		Content:   row.Content,
		CreatedAt: row.CreatedAt,
	}, nil
}

func (r *PostgresArticleInteractionsRepository) DeleteComment(commentID string) error {
	return r.db.Exec(`DELETE FROM article_comments WHERE id = ?::uuid`, commentID).Error
}
