package database

import (
	"fmt"
	"log"

	"stocet-core-service/internal/config"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

// InitDB initializes the PostgreSQL database connection.
// On subsequent calls (e.g. Lambda warm starts) the existing connection is reused.
func InitDB() error {
	if DB != nil {
		return nil
	}

	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=%s",
		config.AppConfig.PostgresHost,
		config.AppConfig.PostgresUser,
		config.AppConfig.PostgresPassword,
		config.AppConfig.PostgresDB,
		config.AppConfig.PostgresPort,
		config.AppConfig.PostgresSSLMode,
	)

	logLevel := logger.Info
	if config.AppConfig.Environment == "production" {
		logLevel = logger.Warn
	}

	var err error
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logLevel),
	})

	if err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}

	log.Println("Connected to PostgreSQL database")

	// NOTE: Schema migrations are managed by golang-migrate.
	// Run `make migrate-up` to apply pending migrations.

	return nil
}

// CloseDB closes the database connection
func CloseDB() error {
	sqlDB, err := DB.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}
