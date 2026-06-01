package repositories

import (
	"time"

	"stocet-core-service/internal/database"
	"stocet-core-service/internal/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type AIQuotaRepository struct{}

func NewAIQuotaRepository() *AIQuotaRepository {
	return &AIQuotaRepository{}
}

func todayUTC() time.Time {
	now := time.Now().UTC()
	y, m, d := now.Date()
	return time.Date(y, m, d, 0, 0, 0, 0, time.UTC)
}

// TryIncrement increments today's count when below limit. Returns new count and whether increment succeeded.
func (r *AIQuotaRepository) TryIncrement(userID string, limit int) (int, bool, error) {
	uid, err := uuid.Parse(userID)
	if err != nil {
		return 0, false, err
	}
	today := todayUTC()

	var newCount int
	var allowed bool
	err = database.DB.Transaction(func(tx *gorm.DB) error {
		var row models.AIUsageDaily
		err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
			Where("user_id = ? AND usage_date = ?", uid, today).
			First(&row).Error
		if err == gorm.ErrRecordNotFound {
			row = models.AIUsageDaily{
				UserID:       uid,
				UsageDate:    today,
				RequestCount: 1,
			}
			if err := tx.Create(&row).Error; err != nil {
				return err
			}
			newCount = 1
			allowed = true
			return nil
		}
		if err != nil {
			return err
		}
		newCount = row.RequestCount
		if row.RequestCount >= limit {
			allowed = false
			return nil
		}
		row.RequestCount++
		newCount = row.RequestCount
		allowed = true
		return tx.Save(&row).Error
	})
	if err != nil {
		return 0, false, err
	}
	return newCount, allowed, nil
}

func (r *AIQuotaRepository) GetCount(userID string) (int, error) {
	uid, err := uuid.Parse(userID)
	if err != nil {
		return 0, err
	}
	today := todayUTC()

	var row models.AIUsageDaily
	err = database.DB.Where("user_id = ? AND usage_date = ?", uid, today).First(&row).Error
	if err == gorm.ErrRecordNotFound {
		return 0, nil
	}
	if err != nil {
		return 0, err
	}
	return row.RequestCount, nil
}
