package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	// Server
	ServerPort string
	ServerHost string

	// Database
	PostgresHost     string
	PostgresPort     string
	PostgresUser     string
	PostgresPassword string
	PostgresDB       string
	PostgresSSLMode  string

	// MongoDB (for reading news/market data)
	MongoDBURL    string
	MongoDBDBName string

	// JWT
	JWTSecret string

	// Environment
	Environment string
}

var AppConfig *Config

func LoadConfig() error {
	// Load .env file if it exists
	_ = godotenv.Load()

	AppConfig = &Config{
		ServerPort: getEnv("SERVER_PORT", "8080"),
		ServerHost: getEnv("SERVER_HOST", "0.0.0.0"),

		PostgresHost:     getEnv("POSTGRES_HOST", "localhost"),
		PostgresPort:     getEnv("POSTGRES_PORT", "5432"),
		PostgresUser:     getEnv("POSTGRES_USER", "postgres"),
		PostgresPassword: getEnv("POSTGRES_PASSWORD", "postgres"),
		PostgresDB:       getEnv("POSTGRES_DB", "stocet_core"),
		PostgresSSLMode:  getEnv("POSTGRES_SSLMODE", "disable"),

		MongoDBURL:    getEnv("MONGODB_URL", "mongodb://localhost:27017"),
		MongoDBDBName: getEnv("MONGODB_DB_NAME", "ethiopian_financial_news"),

		JWTSecret:  getEnv("JWT_SECRET", "change-me-in-production"),
		Environment: getEnv("ENVIRONMENT", "development"),
	}

	if AppConfig.JWTSecret == "change-me-in-production" && AppConfig.Environment == "production" {
		log.Fatal("JWT_SECRET must be set in production")
	}

	return nil
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
