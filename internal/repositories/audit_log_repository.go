package repositories

import (
	"stocet-core-service/internal/database"
	"stocet-core-service/internal/models"
)

const maxAuditLogs = 5000

type AuditLogFilters struct {
	ActorID string
	Method  string
	Limit   int
	Skip    int
}

type AuditLogRepository struct{}

func NewAuditLogRepository() *AuditLogRepository { return &AuditLogRepository{} }

// Create persists one audit entry. Called from a goroutine — errors are silently dropped.
func (r *AuditLogRepository) Create(entry *models.AuditLog) error {
	if err := database.DB.Create(entry).Error; err != nil {
		return err
	}
	// Rolling cap: keep only the most recent maxAuditLogs rows.
	_ = database.DB.Exec(`
		DELETE FROM audit_logs
		WHERE id NOT IN (
			SELECT id FROM audit_logs ORDER BY created_at DESC LIMIT ?
		)
	`, maxAuditLogs)
	return nil
}

// List returns audit entries newest-first with optional filters.
func (r *AuditLogRepository) List(f AuditLogFilters) ([]models.AuditLog, int64, error) {
	if f.Limit <= 0 {
		f.Limit = 50
	}
	if f.Limit > 200 {
		f.Limit = 200
	}

	q := database.DB.Model(&models.AuditLog{})
	if f.ActorID != "" {
		q = q.Where("actor_id = ?", f.ActorID)
	}
	if f.Method != "" {
		q = q.Where("method = ?", f.Method)
	}

	var total int64
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	var logs []models.AuditLog
	err := q.Order("created_at DESC").Limit(f.Limit).Offset(f.Skip).Find(&logs).Error
	return logs, total, err
}
