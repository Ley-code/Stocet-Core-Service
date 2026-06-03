package models

import (
	"time"

	"github.com/google/uuid"
)

type OrgPlan string

const (
	OrgPlanStarter      OrgPlan = "starter"
	OrgPlanProfessional OrgPlan = "professional"
	OrgPlanEnterprise   OrgPlan = "enterprise"
)

type OrgMemberRole string

const (
	OrgRoleOwner  OrgMemberRole = "owner"
	OrgRoleMember OrgMemberRole = "member"
)

type Organization struct {
	ID           uuid.UUID  `gorm:"primaryKey;type:uuid;default:gen_random_uuid()" json:"id"`
	Name         string     `gorm:"type:text;not null"                              json:"name"`
	Slug         string     `gorm:"uniqueIndex;type:varchar(63);not null"           json:"slug"`
	Plan         OrgPlan    `gorm:"type:varchar(20);not null;default:'starter'"     json:"plan"`
	SeatLimit    int        `gorm:"not null;default:10"                             json:"seat_limit"`
	LogoURL      *string    `gorm:"type:text"                                       json:"logo_url,omitempty"`
	BrandColor   *string    `gorm:"type:varchar(7)"                                 json:"brand_color,omitempty"`
	CustomDomain *string    `gorm:"type:varchar(253)"                               json:"custom_domain,omitempty"`
	CreatedBy    *uuid.UUID `gorm:"type:uuid"                                       json:"created_by,omitempty"`
	CreatedAt    time.Time  `                                                       json:"created_at"`
	UpdatedAt    time.Time  `                                                       json:"updated_at"`
}

func (Organization) TableName() string { return "organizations" }

type OrgMember struct {
	ID        uuid.UUID     `gorm:"primaryKey;type:uuid;default:gen_random_uuid()" json:"id"`
	OrgID     uuid.UUID     `gorm:"type:uuid;not null"                             json:"org_id"`
	UserID    uuid.UUID     `gorm:"type:uuid;not null"                             json:"user_id"`
	Role      OrgMemberRole `gorm:"type:varchar(10);not null;default:'member'"     json:"role"`
	CreatedAt time.Time     `                                                      json:"created_at"`
}

func (OrgMember) TableName() string { return "org_members" }
