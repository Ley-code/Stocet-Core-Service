package repositories

import (
	"stocet-core-service/internal/database"
	"stocet-core-service/internal/models"

	"github.com/google/uuid"
)

type NotificationPreferenceRepository struct{}

func NewNotificationPreferenceRepository() *NotificationPreferenceRepository {
	return &NotificationPreferenceRepository{}
}

// Get returns the preferences for a user, creating a default row if none exists.
func (r *NotificationPreferenceRepository) Get(userID string) (*models.NotificationPreference, error) {
	uid, err := uuid.Parse(userID)
	if err != nil {
		return nil, err
	}

	pref := &models.NotificationPreference{
		UserID:          uid,
		PriceAlertPush:  true,
		PriceAlertEmail: false,
		NewsAlertPush:   true,
		NewsAlertEmail:  false,
		InAppEnabled:    true,
	}

	// Upsert: insert defaults if the row doesn't exist yet.
	err = database.DB.Where("user_id = ?", uid).FirstOrCreate(pref).Error
	return pref, err
}

// Save upserts the preferences for a user.
func (r *NotificationPreferenceRepository) Save(pref *models.NotificationPreference) error {
	return database.DB.Save(pref).Error
}
