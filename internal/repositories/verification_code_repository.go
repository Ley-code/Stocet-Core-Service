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

func (r *VerificationCodeRepository) UpsertActiveCode(email, codeHash string, expiresAt time.Time, purpose string) error {
	if purpose == "" {
		purpose = models.VerificationPurposeEmail
	}

	if err := database.DB.Where("email = ? AND purpose = ? AND used_at IS NULL", email, purpose).
		Delete(&models.VerificationCode{}).Error; err != nil {
		return err
	}

	code := &models.VerificationCode{
		Email:     email,
		Purpose:   purpose,
		CodeHash:  codeHash,
		ExpiresAt: expiresAt,
	}

	return database.DB.Create(code).Error
}

func (r *VerificationCodeRepository) GetLatestActiveByEmailAndPurpose(email, purpose string) (*models.VerificationCode, error) {
	if purpose == "" {
		purpose = models.VerificationPurposeEmail
	}

	var code models.VerificationCode
	err := database.DB.
		Where("email = ? AND purpose = ? AND used_at IS NULL", email, purpose).
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
