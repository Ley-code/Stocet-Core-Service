package services

import (
	"stocet-core-service/internal/config"
)

type SubscriptionService struct{}

func NewSubscriptionService() *SubscriptionService {
	return &SubscriptionService{}
}

// HasFeature checks if a user's subscription includes a specific feature
func (s *SubscriptionService) HasFeature(subscriptionStatus, feature string) bool {
	// Free users have basic features
	if subscriptionStatus == config.SUBSCRIPTION_FREE {
		return feature != config.FEATURE_ADVANCED_ANALYTICS &&
			feature != config.FEATURE_PREMIUM_ALERTS &&
			feature != config.FEATURE_EXPORT_DATA &&
			feature != config.FEATURE_API_ACCESS
	}

	// Paid users have most features
	if subscriptionStatus == config.SUBSCRIPTION_PAID {
		return feature != config.FEATURE_API_ACCESS
	}

	// Premium users have all features
	if subscriptionStatus == config.SUBSCRIPTION_PREMIUM {
		return true
	}

	return false
}

// CanAccessFeature is a helper to check feature access
func (s *SubscriptionService) CanAccessFeature(subscriptionStatus, feature string) bool {
	return s.HasFeature(subscriptionStatus, feature)
}
