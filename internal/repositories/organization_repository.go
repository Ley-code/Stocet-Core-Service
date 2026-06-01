package repositories

import (
	"errors"
	"strings"

	"stocet-core-service/internal/database"
	"stocet-core-service/internal/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type OrganizationRepository struct{}

func NewOrganizationRepository() *OrganizationRepository {
	return &OrganizationRepository{}
}

func (r *OrganizationRepository) Create(org *models.Organization) error {
	return database.DB.Create(org).Error
}

func (r *OrganizationRepository) GetByID(id string) (*models.Organization, error) {
	uid, err := uuid.Parse(id)
	if err != nil {
		return nil, errors.New("invalid organization id")
	}
	var org models.Organization
	if err := database.DB.Where("id = ?", uid).First(&org).Error; err != nil {
		return nil, err
	}
	return &org, nil
}

func (r *OrganizationRepository) GetBySlug(slug string) (*models.Organization, error) {
	var org models.Organization
	if err := database.DB.Where("slug = ?", strings.TrimSpace(slug)).First(&org).Error; err != nil {
		return nil, err
	}
	return &org, nil
}

func (r *OrganizationRepository) List() ([]models.Organization, error) {
	var orgs []models.Organization
	err := database.DB.Order("created_at DESC").Find(&orgs).Error
	return orgs, err
}

func (r *OrganizationRepository) Update(org *models.Organization) error {
	return database.DB.Save(org).Error
}

func (r *OrganizationRepository) Delete(id string) error {
	uid, err := uuid.Parse(id)
	if err != nil {
		return errors.New("invalid organization id")
	}
	res := database.DB.Where("id = ?", uid).Delete(&models.Organization{})
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return errors.New("organization not found")
	}
	return nil
}

func (r *OrganizationRepository) CountMembers(orgID uuid.UUID) (int64, error) {
	var count int64
	err := database.DB.Model(&models.OrgMember{}).Where("org_id = ?", orgID).Count(&count).Error
	return count, err
}

func (r *OrganizationRepository) AddMember(orgID, userID uuid.UUID, role models.OrgMemberRole) error {
	member := &models.OrgMember{
		OrgID:  orgID,
		UserID: userID,
		Role:   role,
	}
	return database.DB.Create(member).Error
}

func (r *OrganizationRepository) ListMembers(orgID uuid.UUID) ([]models.OrgMember, error) {
	var members []models.OrgMember
	err := database.DB.Where("org_id = ?", orgID).Order("created_at ASC").Find(&members).Error
	return members, err
}

func (r *OrganizationRepository) RemoveMember(orgID, userID uuid.UUID) error {
	res := database.DB.Where("org_id = ? AND user_id = ?", orgID, userID).Delete(&models.OrgMember{})
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return errors.New("member not found")
	}
	return nil
}

// GetOrgIDByCode returns the org_id stamped on an invite code, if any.
func (r *OrganizationRepository) GetOrgIDByCode(code string) (*uuid.UUID, error) {
	c := strings.TrimSpace(code)
	if c == "" {
		return nil, nil
	}
	var ic models.InviteCode
	err := database.DB.Select("org_id").Where("code = ?", c).First(&ic).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return ic.OrgID, nil
}

// ListInviteCodesByOrg returns invite codes linked to an organization.
func (r *OrganizationRepository) ListInviteCodesByOrg(orgID uuid.UUID) ([]models.InviteCode, error) {
	var codes []models.InviteCode
	err := database.DB.Where("org_id = ?", orgID).Order("created_at DESC").Find(&codes).Error
	return codes, err
}
