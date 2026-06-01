package services

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"math/big"
	"strings"
	"time"

	"stocet-core-service/internal/config"
	"stocet-core-service/internal/email"
	"stocet-core-service/internal/models"
	"stocet-core-service/internal/repositories"
	"stocet-core-service/internal/utils"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// ErrUserNotVerified is returned when a user tries to login but hasn't verified their email.
var ErrUserNotVerified = errors.New("USER_NOT_VERIFIED")

type AuthService struct {
	userRepo             *repositories.UserRepository
	verificationCodeRepo *repositories.VerificationCodeRepository
	inviteCodeRepo       *repositories.InviteCodeRepository
	orgRepo              *repositories.OrganizationRepository
	subSvc               *SubscriptionService
	mailer               email.Sender
}

func NewAuthService(
	userRepo *repositories.UserRepository,
	verificationCodeRepo *repositories.VerificationCodeRepository,
	inviteCodeRepo *repositories.InviteCodeRepository,
	orgRepo *repositories.OrganizationRepository,
	subSvc *SubscriptionService,
	mailer email.Sender,
) *AuthService {
	return &AuthService{
		userRepo:             userRepo,
		verificationCodeRepo: verificationCodeRepo,
		inviteCodeRepo:       inviteCodeRepo,
		orgRepo:              orgRepo,
		subSvc:               subSvc,
		mailer:               mailer,
	}
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type SignupRequest struct {
	Email       string `json:"email" binding:"required"`
	Password    string `json:"password" binding:"required"`
	FullName    string `json:"full_name" binding:"required"`
	PhoneNumber string `json:"phone_number" binding:"required"`
	InviteCode  string `json:"invite_code" binding:"required"`
}

type AuthResponse struct {
	Token string        `json:"token"`
	User  *UserProfile  `json:"user"`
}

type SignupResponse struct {
	Email                string `json:"email"`
	RequiresVerification bool   `json:"requires_verification"`
}

type VerifyEmailRequest struct {
	Email string `json:"email" binding:"required"`
	Code  string `json:"code" binding:"required"`
}

type VerifyEmailResponse struct {
	Email      string `json:"email"`
	IsVerified bool   `json:"is_verified"`
}

func (s *AuthService) Login(req LoginRequest) (*AuthResponse, error) {
	email := strings.ToLower(strings.TrimSpace(req.Email))

	// Validate email
	if !utils.ValidateEmail(email) {
		return nil, errors.New("invalid email format")
	}

	// Get user by email
	user, err := s.userRepo.GetByEmail(email)
	if err != nil {
		return nil, errors.New("invalid credentials")
	}

	if !user.IsVerified {
		return nil, fmt.Errorf("%w:%s", ErrUserNotVerified, email)
	}

	// Check password
	if !utils.CheckPasswordHash(req.Password, user.PasswordHash) {
		return nil, errors.New("invalid credentials")
	}

	effective := string(user.SubscriptionStatus)
	if s.subSvc != nil {
		effective, err = s.subSvc.EffectiveTier(user.ID)
		if err != nil {
			return nil, errors.New("failed to resolve subscription")
		}
	}

	orgID := ""
	if user.OrgID != nil {
		orgID = user.OrgID.String()
	}

	token, err := utils.GenerateToken(user.ID, user.Email, string(user.Role), effective, orgID)
	if err != nil {
		return nil, errors.New("failed to generate token")
	}

	var profile *UserProfile
	if s.subSvc != nil {
		profile, err = s.subSvc.BuildUserProfile(user)
		if err != nil {
			return nil, errors.New("failed to build user profile")
		}
	} else {
		profile = &UserProfile{
			ID:                    user.ID,
			Email:                 user.Email,
			FullName:              user.FullName,
			PhoneNumber:           user.PhoneNumber,
			IsVerified:            user.IsVerified,
			Role:                  string(user.Role),
			SubscriptionStatus:    string(user.SubscriptionStatus),
			EffectiveSubscription: effective,
			HasSeenTour:           user.HasSeenTour,
			CreatedAt:             user.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
			UpdatedAt:             user.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
		}
		if user.OrgID != nil {
			s := user.OrgID.String()
			profile.OrgID = &s
		}
	}

	return &AuthResponse{
		Token: token,
		User:  profile,
	}, nil
}

func (s *AuthService) Signup(req SignupRequest) (*SignupResponse, error) {
	email := strings.ToLower(strings.TrimSpace(req.Email))
	fullName := strings.TrimSpace(req.FullName)
	phoneNumber := strings.TrimSpace(req.PhoneNumber)

	// Validate email
	if !utils.ValidateEmail(email) {
		return nil, errors.New("invalid email format")
	}

	// Validate password
	if !utils.ValidatePassword(req.Password) {
		return nil, errors.New("password must be at least 8 characters with at least one letter and one number")
	}

	if fullName == "" {
		return nil, errors.New("full_name is required")
	}

	if phoneNumber == "" {
		return nil, errors.New("phone_number is required")
	}

	// Check if user already exists
	existingUser, err := s.userRepo.GetByEmail(email)
	if err == nil && existingUser != nil {
		return nil, errors.New("user with this email already exists")
	}
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, errors.New("failed to validate existing user")
	}

	inviteTrim := strings.TrimSpace(req.InviteCode)
	if err := s.inviteCodeRepo.Consume(inviteTrim); err != nil {
		return nil, err
	}

	orgID, _ := s.orgRepo.GetOrgIDByCode(inviteTrim)

	// Hash password
	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		_ = s.inviteCodeRepo.Release(inviteTrim)
		return nil, errors.New("failed to hash password")
	}

	// Create user
	user := &models.User{
		ID:                 uuid.New().String(),
		Email:              email,
		FullName:           fullName,
		PhoneNumber:        phoneNumber,
		IsVerified:         false,
		PasswordHash:       hashedPassword,
		Role:               models.UserRole(config.USER_ROLE),
		SubscriptionStatus: models.SubscriptionStatus(config.SUBSCRIPTION_FREE),
	}

	if err := s.userRepo.Create(user); err != nil {
		_ = s.inviteCodeRepo.Release(inviteTrim)
		return nil, errors.New("failed to create user")
	}

	code, err := generateVerificationCode()
	if err != nil {
		_ = s.userRepo.Delete(user.ID)
		_ = s.inviteCodeRepo.Release(inviteTrim)
		return nil, errors.New("failed to generate verification code")
	}

	codeHash := hashVerificationCode(code)
	expiresAt := time.Now().UTC().Add(3 * time.Minute)

	if err := s.verificationCodeRepo.UpsertActiveCode(email, codeHash, expiresAt); err != nil {
		_ = s.userRepo.Delete(user.ID)
		_ = s.inviteCodeRepo.Release(inviteTrim)
		return nil, errors.New("failed to store verification code")
	}

	if err := s.sendVerificationEmail(email, fullName, code); err != nil {
		_ = s.userRepo.Delete(user.ID)
		_ = s.inviteCodeRepo.Release(inviteTrim)
		return nil, errors.New("failed to send verification email")
	}

	if orgID != nil {
		userUUID, parseErr := uuid.Parse(user.ID)
		if parseErr == nil {
			_ = s.orgRepo.AddMember(*orgID, userUUID, models.OrgRoleMember)
			_ = s.userRepo.SetOrgID(user.ID, *orgID)
		}
	}

	return &SignupResponse{
		Email:                email,
		RequiresVerification: true,
	}, nil
}

func (s *AuthService) VerifyEmail(req VerifyEmailRequest) (*VerifyEmailResponse, error) {
	email := strings.ToLower(strings.TrimSpace(req.Email))
	code := strings.TrimSpace(req.Code)

	if !utils.ValidateEmail(email) {
		return nil, errors.New("invalid email format")
	}
	if len(code) != 6 {
		return nil, errors.New("verification code must be 6 digits")
	}

	verificationCode, err := s.verificationCodeRepo.GetLatestActiveByEmail(email)
	if err != nil {
		return nil, errors.New("invalid or expired verification code")
	}

	if verificationCode.ExpiresAt.Before(time.Now().UTC()) {
		return nil, errors.New("verification code has expired")
	}

	if verificationCode.CodeHash != hashVerificationCode(code) {
		return nil, errors.New("invalid or expired verification code")
	}

	if err := s.userRepo.SetVerifiedByEmail(email); err != nil {
		return nil, errors.New("failed to verify email")
	}

	if err := s.verificationCodeRepo.MarkUsed(verificationCode.ID); err != nil {
		return nil, errors.New("failed to finalize verification")
	}

	return &VerifyEmailResponse{
		Email:      email,
		IsVerified: true,
	}, nil
}

func generateVerificationCode() (string, error) {
	max := big.NewInt(1000000)
	n, err := rand.Int(rand.Reader, max)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%06d", n.Int64()), nil
}

func hashVerificationCode(code string) string {
	sum := sha256.Sum256([]byte(code))
	return hex.EncodeToString(sum[:])
}

func (s *AuthService) sendVerificationEmail(toEmail, fullName, code string) error {
	if s.mailer == nil || !s.mailer.Enabled() {
		return errors.New("email delivery is currently unavailable")
	}

	subject := "Verify your STOCET account"
	htmlBody := fmt.Sprintf(`
<!doctype html>
<html>
  <body style="margin:0;padding:0;background:#f6f8fb;font-family:Arial,sans-serif;color:#111827;">
    <table role="presentation" width="100%%" cellpadding="0" cellspacing="0" style="padding:24px 12px;">
      <tr>
        <td align="center">
          <table role="presentation" width="560" cellpadding="0" cellspacing="0" style="background:#ffffff;border:1px solid #e5e7eb;border-radius:10px;padding:28px;">
            <tr><td style="font-size:20px;font-weight:700;color:#111827;padding-bottom:8px;">STOCET Account Verification</td></tr>
            <tr><td style="font-size:14px;color:#374151;padding-bottom:20px;">Hi %s, thanks for signing up. Use the verification code below to activate your account.</td></tr>
            <tr>
              <td align="center" style="padding:8px 0 20px 0;">
                <div style="display:inline-block;font-size:28px;letter-spacing:6px;font-weight:700;padding:12px 20px;border:1px solid #d1d5db;border-radius:8px;background:#f9fafb;">%s</div>
              </td>
            </tr>
            <tr><td style="font-size:13px;color:#6b7280;padding-bottom:8px;">This code expires in 3 minutes.</td></tr>
            <tr><td style="font-size:12px;color:#9ca3af;">If you did not create this account, you can safely ignore this email.</td></tr>
          </table>
        </td>
      </tr>
    </table>
  </body>
</html>`, fullName, code)

	return s.mailer.SendHTML(toEmail, subject, htmlBody)
}

// ResendVerification generates a new verification code and sends it to the user's email.
func (s *AuthService) ResendVerification(emailAddr string) error {
	emailAddr = strings.ToLower(strings.TrimSpace(emailAddr))

	if !utils.ValidateEmail(emailAddr) {
		return errors.New("invalid email format")
	}

	user, err := s.userRepo.GetByEmail(emailAddr)
	if err != nil {
		return errors.New("no account found with this email address")
	}

	if user.IsVerified {
		return errors.New("this email is already verified")
	}

	code, err := generateVerificationCode()
	if err != nil {
		return errors.New("failed to generate verification code")
	}

	codeHash := hashVerificationCode(code)
	expiresAt := time.Now().UTC().Add(3 * time.Minute)

	if err := s.verificationCodeRepo.UpsertActiveCode(emailAddr, codeHash, expiresAt); err != nil {
		return errors.New("failed to store verification code")
	}

	if err := s.sendVerificationEmail(emailAddr, user.DisplayName(), code); err != nil {
		return errors.New("failed to send verification email")
	}

	return nil
}
