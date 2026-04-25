package models

import (
	"time"

	"github.com/lib/pq"
	"gorm.io/gorm"
)

type CapitalServiceProvider struct {
	ID               string         `gorm:"primaryKey;type:uuid;default:gen_random_uuid()" json:"id"`
	Name             string         `gorm:"not null;index" json:"name"`
	LicenseNumber    string         `gorm:"unique;not null" json:"license_number"`
	Category         string         `gorm:"not null" json:"category"`
	Status           string         `gorm:"default:Active" json:"status"`
	RegistrationDate time.Time      `json:"registration_date"`
	Services         pq.StringArray `gorm:"type:text[]" json:"services"`
	PaidUpCapital    float64        `json:"paid_up_capital"`
	ChiefExecutive   string         `json:"chief_executive"`
	Phone            string         `json:"phone"`
	Email            string         `json:"email"`
	Website          string         `json:"website"`
	TelegramHandle   string         `json:"telegram_handle"`
	HeadOffice       string         `json:"head_office"`
	LogoURL          string         `json:"logo_url"`
	Description      string         `gorm:"type:text" json:"description"`
	CreatedAt        time.Time      `json:"created_at"`
	UpdatedAt        time.Time      `json:"updated_at"`
	DeletedAt        gorm.DeletedAt `gorm:"index" json:"-"`
}

func (CapitalServiceProvider) TableName() string {
	return "capital_service_providers"
}
