package repositories

import (
	"errors"

	"stocet-core-service/internal/database"
	"stocet-core-service/internal/models"

	"github.com/google/uuid"
)

var ErrAdvisoryReportNotFound = errors.New("advisory report not found")

const maxAdvisoryReportsPerUser = 20

type AdvisoryReportRepository struct{}

func NewAdvisoryReportRepository() *AdvisoryReportRepository {
	return &AdvisoryReportRepository{}
}

// Create persists a new report. Enforces a per-user cap by deleting the oldest
// when the user already has maxAdvisoryReportsPerUser reports.
func (r *AdvisoryReportRepository) Create(report *models.AdvisoryReport) error {
	if err := database.DB.Create(report).Error; err != nil {
		return err
	}
	// Evict oldest beyond cap (best-effort; ignore errors)
	_ = database.DB.Exec(`
		DELETE FROM advisory_reports
		WHERE user_id = ? AND id NOT IN (
			SELECT id FROM advisory_reports
			WHERE user_id = ?
			ORDER BY created_at DESC
			LIMIT ?
		)
	`, report.UserID, report.UserID, maxAdvisoryReportsPerUser)
	return nil
}

// ListByUser returns up to 20 reports for a user, newest first.
func (r *AdvisoryReportRepository) ListByUser(userID string) ([]models.AdvisoryReport, error) {
	uid, err := uuid.Parse(userID)
	if err != nil {
		return nil, errors.New("invalid user id")
	}
	var reports []models.AdvisoryReport
	err = database.DB.
		Where("user_id = ?", uid).
		Order("created_at DESC").
		Limit(maxAdvisoryReportsPerUser).
		Find(&reports).Error
	return reports, err
}

// GetByID returns a single report belonging to the given user.
func (r *AdvisoryReportRepository) GetByID(userID, reportID string) (*models.AdvisoryReport, error) {
	uid, err := uuid.Parse(userID)
	if err != nil {
		return nil, errors.New("invalid user id")
	}
	rid, err := uuid.Parse(reportID)
	if err != nil {
		return nil, errors.New("invalid report id")
	}
	var report models.AdvisoryReport
	err = database.DB.Where("id = ? AND user_id = ?", rid, uid).First(&report).Error
	if err != nil {
		return nil, ErrAdvisoryReportNotFound
	}
	return &report, nil
}

// Delete permanently removes a report belonging to the given user.
func (r *AdvisoryReportRepository) Delete(userID, reportID string) error {
	uid, err := uuid.Parse(userID)
	if err != nil {
		return errors.New("invalid user id")
	}
	rid, err := uuid.Parse(reportID)
	if err != nil {
		return errors.New("invalid report id")
	}
	res := database.DB.
		Where("id = ? AND user_id = ?", rid, uid).
		Delete(&models.AdvisoryReport{})
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return ErrAdvisoryReportNotFound
	}
	return nil
}
