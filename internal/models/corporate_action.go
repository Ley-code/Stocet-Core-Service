package models

import (
	"time"

	"github.com/google/uuid"
)

type CorporateActionType string

const (
	CorporateActionAGM         CorporateActionType = "agm"
	CorporateActionDividend    CorporateActionType = "dividend"
	CorporateActionRightsIssue CorporateActionType = "rights_issue"
	CorporateActionBonusShares CorporateActionType = "bonus_shares"
	CorporateActionListing     CorporateActionType = "listing"
	CorporateActionOther       CorporateActionType = "other"
)

type CorporateAction struct {
	ID               uuid.UUID           `gorm:"primaryKey;type:uuid;default:gen_random_uuid()" json:"id"`
	CompanyProfileID string              `gorm:"type:uuid;not null;index"                       json:"company_profile_id"`
	ActionType       CorporateActionType `gorm:"type:corporate_action_type;not null"            json:"action_type"`
	Title            string              `gorm:"type:text;not null"                             json:"title"`
	EventDate        time.Time           `gorm:"type:date;not null"                             json:"event_date"`
	ExDate           *time.Time          `gorm:"type:date"                                      json:"ex_date,omitempty"`
	RecordDate       *time.Time          `gorm:"type:date"                                      json:"record_date,omitempty"`
	Description      *string             `gorm:"type:text"                                      json:"description,omitempty"`
	Amount           *float64            `gorm:"type:numeric(24,6)"                             json:"amount,omitempty"`
	Currency         string              `gorm:"type:varchar(10);not null;default:'ETB'"        json:"currency"`
	// Joined from company_profiles — populated by repository, not stored
	TickerSymbol *string `gorm:"-" json:"ticker_symbol,omitempty"`
	CompanyName  *string `gorm:"-" json:"company_name,omitempty"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

func (CorporateAction) TableName() string { return "corporate_actions" }
