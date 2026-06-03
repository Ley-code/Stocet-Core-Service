package database

import (
	"errors"
	"log"
	"os"
	"stocet-core-service/internal/config"
	"stocet-core-service/internal/models"
	"stocet-core-service/internal/utils"
)

// SeedInitialAdmin creates an initial admin user when ENABLE_ADMIN_SEED=true.
func SeedInitialAdmin() error {
	if !config.AppConfig.EnableAdminSeed {
		log.Println("Admin seed skipped (ENABLE_ADMIN_SEED is not true)")
		return nil
	}

	password := config.AppConfig.AdminInitialPassword
	if password == "" {
		password = os.Getenv("ADMIN_INITIAL_PASSWORD")
	}
	if password == "" {
		return errors.New("ADMIN_INITIAL_PASSWORD is required when ENABLE_ADMIN_SEED=true")
	}

	var existingAdmin models.User
	err := DB.Where("email = ?", "admin@stocet.com").First(&existingAdmin).Error
	if err == nil {
		log.Println("Admin user already exists")
		return nil
	}

	hashedPassword, err := utils.HashPassword(password)
	if err != nil {
		return err
	}

	admin := &models.User{
		Email:              "admin@stocet.com",
		FullName:           "Admin",
		PhoneNumber:        "",
		IsVerified:         true,
		PasswordHash:       hashedPassword,
		Role:               config.ADMIN_ROLE,
		SubscriptionStatus: config.SUBSCRIPTION_PREMIUM,
	}

	if err := DB.Create(admin).Error; err != nil {
		return err
	}

	log.Println("Initial admin user created for admin@stocet.com")
	return nil
}
