package services

import (
	"errors"
	"regexp"
	"strings"

	"stocet-core-service/internal/models"
	"stocet-core-service/internal/repositories"

	"github.com/google/uuid"
)

var slugPattern = regexp.MustCompile(`^[a-z0-9]+(?:-[a-z0-9]+)*$`)

type OrganizationService struct {
	orgRepo  *repositories.OrganizationRepository
	userRepo *repositories.UserRepository
}

func NewOrganizationService(
	orgRepo *repositories.OrganizationRepository,
	userRepo *repositories.UserRepository,
) *OrganizationService {
	return &OrganizationService{orgRepo: orgRepo, userRepo: userRepo}
}

type CreateOrgRequest struct {
	Name      string `json:"name"`
	Slug      string `json:"slug"`
	Plan      string `json:"plan"`
	SeatLimit int    `json:"seat_limit"`
	LogoURL   string `json:"logo_url"`
}

type UpdateOrgRequest struct {
	Name      *string `json:"name"`
	Slug      *string `json:"slug"`
	Plan      *string `json:"plan"`
	SeatLimit *int    `json:"seat_limit"`
	LogoURL   *string `json:"logo_url"`
}

type OrgListItem struct {
	models.Organization
	MemberCount int64 `json:"member_count"`
}

type OrgWithMembers struct {
	models.Organization
	Members     []OrgMemberWithUser `json:"members"`
	MemberCount int64               `json:"member_count"`
	InviteCodes []models.InviteCode `json:"invite_codes"`
}

type OrgMemberWithUser struct {
	models.OrgMember
	Email    string `json:"email"`
	FullName string `json:"full_name"`
}

func normalizeSlug(slug string) string {
	s := strings.ToLower(strings.TrimSpace(slug))
	s = strings.ReplaceAll(s, " ", "-")
	return s
}

func validateSlug(slug string) error {
	if slug == "" || len(slug) > 63 {
		return errors.New("slug must be 1-63 characters")
	}
	if !slugPattern.MatchString(slug) {
		return errors.New("slug must be lowercase alphanumeric with hyphens")
	}
	return nil
}

func parsePlan(plan string) (models.OrgPlan, error) {
	switch models.OrgPlan(strings.TrimSpace(plan)) {
	case models.OrgPlanStarter, models.OrgPlanProfessional, models.OrgPlanEnterprise:
		return models.OrgPlan(plan), nil
	case "":
		return models.OrgPlanStarter, nil
	default:
		return "", errors.New("plan must be starter, professional, or enterprise")
	}
}

func (s *OrganizationService) CreateOrg(req CreateOrgRequest, createdBy uuid.UUID) (*models.Organization, error) {
	name := strings.TrimSpace(req.Name)
	if name == "" {
		return nil, errors.New("name is required")
	}

	slug := normalizeSlug(req.Slug)
	if slug == "" {
		slug = normalizeSlug(name)
	}
	if err := validateSlug(slug); err != nil {
		return nil, err
	}

	plan, err := parsePlan(req.Plan)
	if err != nil {
		return nil, err
	}

	seatLimit := req.SeatLimit
	if seatLimit <= 0 {
		seatLimit = 10
	}

	var logoURL *string
	if u := strings.TrimSpace(req.LogoURL); u != "" {
		logoURL = &u
	}

	org := &models.Organization{
		Name:      name,
		Slug:      slug,
		Plan:      plan,
		SeatLimit: seatLimit,
		LogoURL:   logoURL,
		CreatedBy: &createdBy,
	}

	if err := s.orgRepo.Create(org); err != nil {
		if strings.Contains(err.Error(), "unique") || strings.Contains(err.Error(), "duplicate") {
			return nil, errors.New("an organization with that slug already exists")
		}
		return nil, err
	}

	return org, nil
}

func (s *OrganizationService) ListOrgs() ([]OrgListItem, error) {
	orgs, err := s.orgRepo.List()
	if err != nil {
		return nil, err
	}

	items := make([]OrgListItem, 0, len(orgs))
	for _, org := range orgs {
		count, err := s.orgRepo.CountMembers(org.ID)
		if err != nil {
			return nil, err
		}
		items = append(items, OrgListItem{Organization: org, MemberCount: count})
	}
	return items, nil
}

func (s *OrganizationService) GetOrgWithMembers(id string) (*OrgWithMembers, error) {
	org, err := s.orgRepo.GetByID(id)
	if err != nil {
		return nil, errors.New("organization not found")
	}

	members, err := s.orgRepo.ListMembers(org.ID)
	if err != nil {
		return nil, err
	}

	count, err := s.orgRepo.CountMembers(org.ID)
	if err != nil {
		return nil, err
	}

	enriched := make([]OrgMemberWithUser, 0, len(members))
	for _, m := range members {
		user, err := s.userRepo.GetByID(m.UserID.String())
		if err != nil {
			continue
		}
		enriched = append(enriched, OrgMemberWithUser{
			OrgMember: m,
			Email:     user.Email,
			FullName:  user.FullName,
		})
	}

	codes, err := s.orgRepo.ListInviteCodesByOrg(org.ID)
	if err != nil {
		return nil, err
	}

	return &OrgWithMembers{
		Organization: *org,
		Members:      enriched,
		MemberCount:  count,
		InviteCodes:  codes,
	}, nil
}

func (s *OrganizationService) UpdateOrg(id string, req UpdateOrgRequest) (*models.Organization, error) {
	org, err := s.orgRepo.GetByID(id)
	if err != nil {
		return nil, errors.New("organization not found")
	}

	if req.Name != nil {
		name := strings.TrimSpace(*req.Name)
		if name == "" {
			return nil, errors.New("name cannot be empty")
		}
		org.Name = name
	}
	if req.Slug != nil {
		slug := normalizeSlug(*req.Slug)
		if err := validateSlug(slug); err != nil {
			return nil, err
		}
		org.Slug = slug
	}
	if req.Plan != nil {
		plan, err := parsePlan(*req.Plan)
		if err != nil {
			return nil, err
		}
		org.Plan = plan
	}
	if req.SeatLimit != nil {
		if *req.SeatLimit <= 0 {
			return nil, errors.New("seat_limit must be positive")
		}
		org.SeatLimit = *req.SeatLimit
	}
	if req.LogoURL != nil {
		u := strings.TrimSpace(*req.LogoURL)
		if u == "" {
			org.LogoURL = nil
		} else {
			org.LogoURL = &u
		}
	}

	if err := s.orgRepo.Update(org); err != nil {
		if strings.Contains(err.Error(), "unique") || strings.Contains(err.Error(), "duplicate") {
			return nil, errors.New("an organization with that slug already exists")
		}
		return nil, err
	}
	return org, nil
}

func (s *OrganizationService) DeleteOrg(id string) error {
	return s.orgRepo.Delete(id)
}

func (s *OrganizationService) AddMember(orgID, userID string, role string) error {
	orgUID, err := uuid.Parse(orgID)
	if err != nil {
		return errors.New("invalid organization id")
	}
	userUID, err := uuid.Parse(userID)
	if err != nil {
		return errors.New("invalid user id")
	}

	org, err := s.orgRepo.GetByID(orgID)
	if err != nil {
		return errors.New("organization not found")
	}

	if _, err := s.userRepo.GetByID(userID); err != nil {
		return errors.New("user not found")
	}

	count, err := s.orgRepo.CountMembers(orgUID)
	if err != nil {
		return err
	}
	if int(count) >= org.SeatLimit {
		return errors.New("organization seat limit reached")
	}

	memberRole := models.OrgRoleMember
	switch models.OrgMemberRole(strings.TrimSpace(role)) {
	case models.OrgRoleOwner:
		memberRole = models.OrgRoleOwner
	case models.OrgRoleMember, "":
		memberRole = models.OrgRoleMember
	default:
		return errors.New("role must be owner or member")
	}

	if err := s.orgRepo.AddMember(orgUID, userUID, memberRole); err != nil {
		if strings.Contains(err.Error(), "unique") || strings.Contains(err.Error(), "duplicate") {
			return errors.New("user is already a member of this organization")
		}
		return err
	}

	return s.userRepo.SetOrgID(userID, orgUID)
}

func (s *OrganizationService) RemoveMember(orgID, userID string) error {
	orgUID, err := uuid.Parse(orgID)
	if err != nil {
		return errors.New("invalid organization id")
	}
	userUID, err := uuid.Parse(userID)
	if err != nil {
		return errors.New("invalid user id")
	}

	if err := s.orgRepo.RemoveMember(orgUID, userUID); err != nil {
		return err
	}

	user, err := s.userRepo.GetByID(userID)
	if err != nil {
		return nil
	}
	if user.OrgID != nil && *user.OrgID == orgUID {
		return s.userRepo.UpdateFields(userID, map[string]interface{}{"org_id": nil})
	}
	return nil
}
