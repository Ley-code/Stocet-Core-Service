package database

import (
	"log"
	"stocet-core-service/internal/config"
	"stocet-core-service/internal/models"
	"stocet-core-service/internal/utils"

	"github.com/google/uuid"
)

// SeedInitialAdmin creates an initial admin user if one doesn't exist
func SeedInitialAdmin() error {
	// Check if admin user exists
	var existingAdmin models.User
	err := DB.Where("email = ?", "admin@stocet.com").First(&existingAdmin).Error
	if err == nil {
		log.Println("Admin user already exists")
		return nil
	}

	// Create admin user
	hashedPassword, err := utils.HashPassword("admin123") // Change this in production!
	if err != nil {
		return err
	}

	admin := &models.User{
		ID:                uuid.New().String(),
		Email:             "admin@stocet.com",
		PasswordHash:      hashedPassword,
		Role:              config.ADMIN_ROLE,
		SubscriptionStatus: config.SUBSCRIPTION_PREMIUM,
	}

	if err := DB.Create(admin).Error; err != nil {
		return err
	}

	log.Println("Initial admin user created: admin@stocet.com / admin123")
	return nil
}
