package repositories

import (
	"errors"
	"strings"

	"stocet-core-service/internal/database"
)

// ErrInvalidInviteCode is returned when an invite cannot be consumed.
var ErrInvalidInviteCode = errors.New("Invalid or expired invite code.")

type InviteCodeRepository struct{}

func NewInviteCodeRepository() *InviteCodeRepository {
	return &InviteCodeRepository{}
}

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
