package controllers

import (
	"errors"
	"net/http"
	"strings"
	"stocet-core-service/internal/repositories"
	"stocet-core-service/internal/services"
	"stocet-core-service/internal/utils"

	"github.com/gin-gonic/gin"
)

type AuthController struct {
	service  *services.AuthService
	userRepo *repositories.UserRepository
	subSvc   *services.SubscriptionService
}

func NewAuthController(
	service *services.AuthService,
	userRepo *repositories.UserRepository,
	subSvc *services.SubscriptionService,
) *AuthController {
	return &AuthController{service: service, userRepo: userRepo, subSvc: subSvc}
}

func (c *AuthController) Login(ctx *gin.Context) {
	var req services.LoginRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(ctx, "Invalid request body")
		return
	}

	response, err := c.service.Login(req)
	if err != nil {
		// Check for unverified user sentinel error
		if errors.Is(err, services.ErrUserNotVerified) {
			// Extract email from the wrapped error "USER_NOT_VERIFIED:<email>"
			parts := strings.SplitN(err.Error(), ":", 2)
			email := ""
			if len(parts) == 2 {
				email = parts[1]
			}
			ctx.JSON(http.StatusForbidden, gin.H{
				"success":    false,
				"error_code": "USER_NOT_VERIFIED",
				"email":      email,
				"message":    "Your email address has not been verified. Please check your inbox for the verification code.",
			})
			return
		}
		utils.Unauthorized(ctx, err.Error())
		return
	}

	utils.SuccessResponse(ctx, 200, response, "Login successful")
}

func (c *AuthController) Signup(ctx *gin.Context) {
	var req services.SignupRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(ctx, "Invalid request body")
		return
	}

	response, err := c.service.Signup(req)
	if err != nil {
		utils.BadRequest(ctx, err.Error())
		return
	}

	utils.SuccessResponse(ctx, 201, response, "Signup successful")
}

func (c *AuthController) VerifyEmail(ctx *gin.Context) {
	var req services.VerifyEmailRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(ctx, "Invalid request body")
		return
	}

	response, err := c.service.VerifyEmail(req)
	if err != nil {
		utils.BadRequest(ctx, err.Error())
		return
	}

	utils.SuccessResponse(ctx, 200, response, "Email verified successfully")
}

type ResendVerificationRequest struct {
	Email string `json:"email" binding:"required"`
}

func (c *AuthController) ResendVerification(ctx *gin.Context) {
	var req ResendVerificationRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(ctx, "Please provide a valid email address")
		return
	}

	if err := c.service.ResendVerification(req.Email); err != nil {
		utils.BadRequest(ctx, err.Error())
		return
	}

	utils.SuccessResponse(ctx, 200, nil, "A new verification code has been sent to your email")
}

// GetMe returns the currently authenticated user's profile
func (c *AuthController) GetMe(ctx *gin.Context) {
	userID, _ := ctx.Get("user_id")

	user, err := c.userRepo.GetByID(userID.(string))
	if err != nil {
		utils.NotFound(ctx, "User not found")
		return
	}

	profile, err := c.subSvc.BuildUserProfile(user)
	if err != nil {
		utils.InternalServerError(ctx, err, "Failed to load profile")
		return
	}

	utils.SuccessResponse(ctx, 200, profile, "")
}

// DeleteAccount permanently deletes the authenticated user's account
func (c *AuthController) DeleteAccount(ctx *gin.Context) {
	userID, _ := ctx.Get("user_id")

	if err := c.userRepo.Delete(userID.(string)); err != nil {
		utils.InternalServerError(ctx, err, "Failed to delete account")
		return
	}

	utils.SuccessResponse(ctx, 200, nil, "Account deleted successfully")
}

// MarkTourSeen marks the onboarding tour as seen for the authenticated user
func (c *AuthController) MarkTourSeen(ctx *gin.Context) {
	userID, _ := ctx.Get("user_id")

	if err := c.userRepo.UpdateFields(userID.(string), map[string]interface{}{
		"has_seen_tour": true,
	}); err != nil {
		utils.InternalServerError(ctx, err, "Failed to update tour status")
		return
	}

	utils.SuccessResponse(ctx, 200, nil, "Tour marked as seen")
}

// UpdateProfile updates the authenticated user's display name
func (c *AuthController) UpdateProfile(ctx *gin.Context) {
	userID, _ := ctx.Get("user_id")

	var req struct {
		FullName string `json:"full_name" binding:"required"`
	}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(ctx, "full_name is required")
		return
	}

	name := strings.TrimSpace(req.FullName)
	if name == "" {
		utils.BadRequest(ctx, "full_name cannot be empty")
		return
	}

	if err := c.userRepo.UpdateFields(userID.(string), map[string]interface{}{
		"full_name": name,
	}); err != nil {
		utils.InternalServerError(ctx, err, "Failed to update profile")
		return
	}

	utils.SuccessResponse(ctx, 200, nil, "Profile updated successfully")
}

// ChangePassword updates the authenticated user's password
func (c *AuthController) ChangePassword(ctx *gin.Context) {
	userID, _ := ctx.Get("user_id")

	var req struct {
		CurrentPassword string `json:"current_password" binding:"required"`
		NewPassword     string `json:"new_password" binding:"required"`
	}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(ctx, "current_password and new_password are required")
		return
	}

	// Get user to verify current password
	user, err := c.userRepo.GetByID(userID.(string))
	if err != nil {
		utils.NotFound(ctx, "User not found")
		return
	}

	if !utils.CheckPasswordHash(req.CurrentPassword, user.PasswordHash) {
		utils.BadRequest(ctx, "Current password is incorrect")
		return
	}

	if !utils.ValidatePassword(req.NewPassword) {
		utils.BadRequest(ctx, "Password must be at least 8 characters with at least one letter and one number")
		return
	}

	hashed, err := utils.HashPassword(req.NewPassword)
	if err != nil {
		utils.InternalServerError(ctx, err, "Failed to hash password")
		return
	}

	if err := c.userRepo.UpdateFields(userID.(string), map[string]interface{}{
		"password_hash": hashed,
	}); err != nil {
		utils.InternalServerError(ctx, err, "Failed to update password")
		return
	}

	utils.SuccessResponse(ctx, 200, nil, "Password changed successfully")
}

// GoogleLogin handles POST /auth/google — verifies a Google ID token and returns a Stocet JWT.
func (c *AuthController) GoogleLogin(ctx *gin.Context) {
	var req services.GoogleLoginRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(ctx, "id_token is required")
		return
	}

	response, err := c.service.GoogleLogin(req)
	if err != nil {
		utils.Unauthorized(ctx, err.Error())
		return
	}

	utils.SuccessResponse(ctx, http.StatusOK, response, "Login successful")
}

// RegisterFCMToken stores the device FCM token for push notifications.
// Called by the mobile app after obtaining a token from Firebase.
func (c *AuthController) RegisterFCMToken(ctx *gin.Context) {
	userID, _ := ctx.Get("user_id")

	var body struct {
		Token string `json:"token" binding:"required"`
	}
	if err := ctx.ShouldBindJSON(&body); err != nil || strings.TrimSpace(body.Token) == "" {
		utils.BadRequest(ctx, "token is required")
		return
	}

	if err := c.userRepo.UpdateFields(userID.(string), map[string]interface{}{
		"fcm_token": body.Token,
	}); err != nil {
		utils.InternalServerError(ctx, err, "Failed to store FCM token")
		return
	}

	utils.SuccessResponse(ctx, http.StatusOK, nil, "FCM token registered")
}
