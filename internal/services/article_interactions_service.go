package services

import (
	"errors"
	"strings"

	"stocet-core-service/internal/repositories"
)

const maxArticleCommentLength = 2000

type ArticleInteractionsService struct {
	repo *repositories.PostgresArticleInteractionsRepository
}

func NewArticleInteractionsService(repo *repositories.PostgresArticleInteractionsRepository) *ArticleInteractionsService {
	return &ArticleInteractionsService{repo: repo}
}

func (s *ArticleInteractionsService) ToggleLike(articleID, userID string) (bool, int64, error) {
	if articleID == "" {
		return false, 0, errors.New("article id is required")
	}

	liked, err := s.repo.ToggleLike(articleID, userID)
	if err != nil {
		return false, 0, err
	}

	summary, err := s.repo.GetLikeStatus(articleID, userID)
	return liked, summary.Count, err
}

func (s *ArticleInteractionsService) GetLikeStatus(articleID, userID string) (bool, int64, error) {
	if articleID == "" {
		return false, 0, errors.New("article id is required")
	}

	summary, err := s.repo.GetLikeStatus(articleID, userID)
	return summary.Liked, summary.Count, err
}

func (s *ArticleInteractionsService) GetCommentCounts(articleIDs []string) (map[string]int64, error) {
	return s.repo.GetCommentCounts(articleIDs)
}

func (s *ArticleInteractionsService) GetLikeCounts(articleIDs []string, userID string) (map[string]repositories.ArticleLikeCountRow, error) {
	return s.repo.GetLikeCounts(articleIDs, userID)
}

func (s *ArticleInteractionsService) CreateComment(articleID, userID, userName, content string) (*repositories.ArticleComment, error) {
	content = strings.TrimSpace(content)
	if articleID == "" {
		return nil, errors.New("article id is required")
	}
	if content == "" {
		return nil, errors.New("content is required")
	}
	if len(content) > maxArticleCommentLength {
		return nil, errors.New("comment too long")
	}

	if userName == "" {
		userName = "User"
	}

	return s.repo.CreateComment(articleID, userID, userName, content)
}

func (s *ArticleInteractionsService) GetComments(articleID string, limit, skip int) ([]repositories.ArticleComment, int64, error) {
	if articleID == "" {
		return nil, 0, errors.New("article id is required")
	}

	return s.repo.GetComments(articleID, limit, skip)
}

func (s *ArticleInteractionsService) DeleteComment(commentID, userID string) error {
	comment, err := s.repo.GetCommentByID(commentID)
	if err != nil {
		return errors.New("comment not found")
	}

	if comment.UserID != userID {
		return errors.New("forbidden")
	}

	return s.repo.DeleteComment(commentID)
}
