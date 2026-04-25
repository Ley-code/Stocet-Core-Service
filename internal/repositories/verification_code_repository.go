package repositories

import (
	"time"

	"stocet-core-service/internal/database"
	"stocet-core-service/internal/models"
)

type VerificationCodeRepository struct{}

func NewVerificationCodeRepository() *VerificationCodeRepository {
	return &VerificationCodeRepository{}
}

func (r *VerificationCodeRepository) UpsertActiveCode(email, codeHash string, expiresAt time.Time) error {
	// Ensure there is only one active OTP per email.
	if err := database.DB.Where("email = ? AND used_at IS NULL", email).Delete(&models.VerificationCode{}).Error; err != nil {
		return err
	}

	code := &models.VerificationCode{
		Email:     email,
		CodeHash:  codeHash,
		ExpiresAt: expiresAt,
	}

	return database.DB.Create(code).Error
}

func (r *VerificationCodeRepository) GetLatestActiveByEmail(email string) (*models.VerificationCode, error) {
	var code models.VerificationCode
	err := database.DB.
		Where("email = ? AND used_at IS NULL", email).
		Order("created_at DESC").
		First(&code).Error
	if err != nil {
		return nil, err
	}
	return &code, nil
}

func (r *VerificationCodeRepository) MarkUsed(id string) error {
	now := time.Now().UTC()
	return database.DB.Model(&models.VerificationCode{}).
		Where("id = ? AND used_at IS NULL", id).
		Update("used_at", now).Error
}
