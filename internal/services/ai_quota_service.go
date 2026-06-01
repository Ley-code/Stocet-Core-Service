package services

import (
	"errors"

	"stocet-core-service/internal/config"
	"stocet-core-service/internal/repositories"
)

var ErrAIQuotaExceeded = errors.New("ai_quota_exceeded")

type AIQuotaService struct {
	repo *repositories.AIQuotaRepository
	sub  *SubscriptionService
}

func NewAIQuotaService(repo *repositories.AIQuotaRepository, sub *SubscriptionService) *AIQuotaService {
	return &AIQuotaService{repo: repo, sub: sub}
}

func dailyLimitForTier(tier string) int {
	switch tier {
	case config.SUBSCRIPTION_FREE:
		return config.AIQuotaFree
	case config.SUBSCRIPTION_PAID:
		return config.AIQuotaPaid
	case config.SUBSCRIPTION_PREMIUM:
		return config.AIQuotaPremium
	default:
		return config.AIQuotaFree
	}
}

// LimitForUser returns the daily AI request cap for the user's effective tier.
func (s *AIQuotaService) LimitForUser(userID string) (int, error) {
	tier, err := s.sub.EffectiveTier(userID)
	if err != nil {
		return 0, err
	}
	return dailyLimitForTier(tier), nil
}

// CheckAndIncrement enforces the daily cap. Admins bypass the limit.
func (s *AIQuotaService) CheckAndIncrement(userID string, isAdmin bool) error {
	if isAdmin {
		return nil
	}
	limit, err := s.LimitForUser(userID)
	if err != nil {
		return err
	}
	if limit <= 0 {
		return ErrAIQuotaExceeded
	}

	_, ok, err := s.repo.TryIncrement(userID, limit)
	if err != nil {
		return err
	}
	if !ok {
		return ErrAIQuotaExceeded
	}
	return nil
}

// GetUsage returns today's used count and daily limit.
func (s *AIQuotaService) GetUsage(userID string, isAdmin bool) (used, limit int, err error) {
	if isAdmin {
		return 0, -1, nil
	}
	limit, err = s.LimitForUser(userID)
	if err != nil {
		return 0, 0, err
	}
	used, err = s.repo.GetCount(userID)
	if err != nil {
		return 0, limit, err
	}
	return used, limit, nil
}
