package database

import (
	"fmt"
	"log"

	"stocet-core-service/internal/config"
	"stocet-core-service/internal/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

// InitDB initializes the PostgreSQL database connection
func InitDB() error {
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=%s",
		config.AppConfig.PostgresHost,
		config.AppConfig.PostgresUser,
		config.AppConfig.PostgresPassword,
		config.AppConfig.PostgresDB,
		config.AppConfig.PostgresPort,
		config.AppConfig.PostgresSSLMode,
	)

	var err error
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})

	if err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}

	log.Println("Connected to PostgreSQL database")

	// Auto-migrate all models
	if err := AutoMigrate(); err != nil {
		return fmt.Errorf("failed to auto-migrate: %w", err)
	}

	return nil
}

// AutoMigrate runs database migrations
func AutoMigrate() error {
	return DB.AutoMigrate(
		&models.User{},
		&models.Broker{},
		&models.Course{},
		&models.Lesson{},
		&models.Quiz{},
		&models.Question{},
		&models.GlossaryTerm{},
		&models.QuizSubmission{},
		&models.Channel{},
		&models.Message{},
		&models.Alert{},
		&models.Subscription{},
	)
}

// CloseDB closes the database connection
func CloseDB() error {
	sqlDB, err := DB.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}
