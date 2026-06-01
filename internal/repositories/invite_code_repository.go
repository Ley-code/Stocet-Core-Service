package repositories

import (
	"crypto/rand"
	"errors"
	"fmt"
	"math/big"
	"strings"

	"stocet-core-service/internal/database"
	"stocet-core-service/internal/models"

	"github.com/google/uuid"
)

var ErrInvalidInviteCode = errors.New("Invalid or expired invite code.")

type InviteCodeRepository struct{}

func NewInviteCodeRepository() *InviteCodeRepository {
	return &InviteCodeRepository{}
}

// ── Public (signup flow) ────────────────────────────────────────────────────

// Consume increments current_uses if the code is active and not exhausted.
func (r *InviteCodeRepository) Consume(code string) error {
	c := strings.TrimSpace(code)
	if c == "" {
		return ErrInvalidInviteCode
	}
	res := database.DB.Exec(`
		UPDATE invite_codes
		SET current_uses = current_uses + 1
		WHERE code = ? AND is_active = true AND current_uses < max_uses
	`, c)
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected != 1 {
		return ErrInvalidInviteCode
	}
	return nil
}

// Release decrements current_uses after a failed signup (best-effort rollback).
func (r *InviteCodeRepository) Release(code string) error {
	c := strings.TrimSpace(code)
	if c == "" {
		return nil
	}
	return database.DB.Exec(`
		UPDATE invite_codes
		SET current_uses = current_uses - 1
		WHERE code = ? AND current_uses > 0
	`, c).Error
}

// ── Admin CRUD ──────────────────────────────────────────────────────────────

// List returns all invite codes ordered by created_at desc.
func (r *InviteCodeRepository) List() ([]models.InviteCode, error) {
	var codes []models.InviteCode
	err := database.DB.Order("created_at DESC").Find(&codes).Error
	return codes, err
}

// Create persists a new invite code. If ic.Code is empty a random code is generated.
func (r *InviteCodeRepository) Create(ic *models.InviteCode) error {
	if strings.TrimSpace(ic.Code) == "" {
		code, err := generateCode()
		if err != nil {
			return err
		}
		ic.Code = code
	}
	ic.Code = strings.ToUpper(strings.TrimSpace(ic.Code))
	return database.DB.Create(ic).Error
}

// Revoke marks an invite code as inactive by its UUID.
func (r *InviteCodeRepository) Revoke(id string) error {
	uid, err := uuid.Parse(id)
	if err != nil {
		return errors.New("invalid invite code id")
	}
	res := database.DB.Model(&models.InviteCode{}).
		Where("id = ?", uid).
		Update("is_active", false)
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return errors.New("invite code not found")
	}
	return nil
}

// Delete permanently removes an invite code by UUID.
func (r *InviteCodeRepository) Delete(id string) error {
	uid, err := uuid.Parse(id)
	if err != nil {
		return errors.New("invalid invite code id")
	}
	res := database.DB.Where("id = ?", uid).Delete(&models.InviteCode{})
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return errors.New("invite code not found")
	}
	return nil
}

// ── Helpers ─────────────────────────────────────────────────────────────────

const codeAlphabet = "ABCDEFGHJKLMNPQRSTUVWXYZ23456789"

func generateCode() (string, error) {
	seg := func(n int) (string, error) {
		b := make([]byte, n)
		for i := range b {
			idx, err := rand.Int(rand.Reader, big.NewInt(int64(len(codeAlphabet))))
			if err != nil {
				return "", err
			}
			b[i] = codeAlphabet[idx.Int64()]
		}
		return string(b), nil
	}
	a, err := seg(5)
	if err != nil {
		return "", err
	}
	b, err := seg(5)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("STOCET-%s-%s", a, b), nil
}
