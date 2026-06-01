package services

import (
	"errors"

	"stocet-core-service/internal/config"
	"stocet-core-service/internal/models"
	"stocet-core-service/internal/repositories"

	"gorm.io/gorm"
)

// ErrSubscriptionRequired is returned when a user's effective tier is below the required level.
var ErrSubscriptionRequired = errors.New("subscription_required")

// ErrAlertLimitReached is returned when a free-tier user has reached the alert cap.
var ErrAlertLimitReached = errors.New("alert_limit_reached")

// subscriptionOrder maps tier strings to numeric levels for comparison.
var subscriptionOrder = map[string]int{
	config.SUBSCRIPTION_FREE:    0,
	config.SUBSCRIPTION_PAID:    1,
	config.SUBSCRIPTION_PREMIUM: 2,
}

const freeTierAlertCap = 3

type SubscriptionService struct {
	userRepo *repositories.UserRepository
	orgRepo  *repositories.OrganizationRepository
}

func NewSubscriptionService(
	userRepo *repositories.UserRepository,
	orgRepo *repositories.OrganizationRepository,
) *SubscriptionService {
	return &SubscriptionService{userRepo: userRepo, orgRepo: orgRepo}
}

// EffectiveTier returns the subscription tier used for feature gating.
// Org members inherit tier from their organization's plan.
func (s *SubscriptionService) EffectiveTier(userID string) (string, error) {
	user, err := s.userRepo.GetByID(userID)
	if err != nil {
		return "", err
	}

	if user.OrgID != nil {
		org, err := s.orgRepo.GetByID(user.OrgID.String())
		if err == nil {
			return orgPlanToTier(org.Plan), nil
		}
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return "", err
		}
	}

	return string(user.SubscriptionStatus), nil
}

func orgPlanToTier(plan models.OrgPlan) string {
	switch plan {
	case models.OrgPlanProfessional, models.OrgPlanEnterprise:
		return config.SUBSCRIPTION_PREMIUM
	default:
		return config.SUBSCRIPTION_FREE
	}
}

// TierLevel returns the numeric level for a tier string.
func TierLevel(tier string) int {
	if lvl, ok := subscriptionOrder[tier]; ok {
		return lvl
	}
	return 0
}

// CanAccess returns true if the user's effective tier meets minLevel.
func (s *SubscriptionService) CanAccess(userID, minLevel string) (bool, string, error) {
	effective, err := s.EffectiveTier(userID)
	if err != nil {
		return false, "", err
	}
	required := TierLevel(minLevel)
	return TierLevel(effective) >= required, effective, nil
}

// MaxAlerts returns the alert cap for a user (-1 = unlimited).
func (s *SubscriptionService) MaxAlerts(userID string) (int, error) {
	effective, err := s.EffectiveTier(userID)
	if err != nil {
		return 0, err
	}
	if TierLevel(effective) >= TierLevel(config.SUBSCRIPTION_PREMIUM) {
		return -1, nil
	}
	return freeTierAlertCap, nil
}

// HasFeature checks if a subscription tier includes a specific feature flag.
func (s *SubscriptionService) HasFeature(subscriptionStatus, feature string) bool {
	if subscriptionStatus == config.SUBSCRIPTION_FREE {
		return feature != config.FEATURE_ADVANCED_ANALYTICS &&
			feature != config.FEATURE_PREMIUM_ALERTS &&
			feature != config.FEATURE_EXPORT_DATA &&
			feature != config.FEATURE_API_ACCESS
	}
	if subscriptionStatus == config.SUBSCRIPTION_PAID {
		return feature != config.FEATURE_API_ACCESS
	}
	if subscriptionStatus == config.SUBSCRIPTION_PREMIUM {
		return true
	}
	return false
}

// UserProfile is the wire shape for /auth/me and login user payloads.
type UserProfile struct {
	ID                    string  `json:"id"`
	Email                 string  `json:"email"`
	FullName              string  `json:"full_name"`
	PhoneNumber           string  `json:"phone_number"`
	IsVerified            bool    `json:"is_verified"`
	Role                  string  `json:"role"`
	SubscriptionStatus    string  `json:"subscription_status"`
	EffectiveSubscription string  `json:"effective_subscription"`
	HasSeenTour           bool    `json:"has_seen_tour"`
	OrgID                 *string `json:"org_id,omitempty"`
	CreatedAt             string  `json:"created_at"`
	UpdatedAt             string  `json:"updated_at"`
}

// BuildUserProfile constructs a profile DTO with computed effective_subscription.
func (s *SubscriptionService) BuildUserProfile(user *models.User) (*UserProfile, error) {
	effective, err := s.EffectiveTier(user.ID)
	if err != nil {
		return nil, err
	}

	profile := &UserProfile{
		ID:                    user.ID,
		Email:                 user.Email,
		FullName:              user.FullName,
		PhoneNumber:           user.PhoneNumber,
		IsVerified:            user.IsVerified,
		Role:                  string(user.Role),
		SubscriptionStatus:    string(user.SubscriptionStatus),
		EffectiveSubscription: effective,
		HasSeenTour:           user.HasSeenTour,
		CreatedAt:             user.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt:             user.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}
	if user.OrgID != nil {
		s := user.OrgID.String()
		profile.OrgID = &s
	}
	return profile, nil
}
