package config

import (
	"log"
	"os"
	"strconv"

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

	// Email (SMTP)
	SMTPHost     string
	SMTPPort     string
	SMTPUsername string
	SMTPPassword string
	SMTPFrom     string

	// AI — Claude (chat / tool-calling)
	ClaudeAPIKey  string
	ClaudeModel   string

	// AI — Gemini (embeddings / RAG only; chat has moved to Claude)
	GeminiAPIKey          string
	GeminiChatModel       string
	GeminiEmbedModel      string
	AIChatHistoryMaxTurns int
	AIChatHistoryMaxChars int
	AIRequestTimeoutSec   int
	AIBackfillIntervalSec int
	AIBackfillBatchSize   int

	// Firebase Cloud Messaging (push notifications — HTTP v1 API)
	FCMProjectID          string
	FCMServiceAccountJSON string
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

		JWTSecret:   getEnv("JWT_SECRET", "change-me-in-production"),
		Environment: getEnv("ENVIRONMENT", "development"),

		SMTPHost:     getEnv("SMTP_HOST", ""),
		SMTPPort:     getEnv("SMTP_PORT", ""),
		SMTPUsername: getEnv("SMTP_USERNAME", ""),
		SMTPPassword: getEnv("SMTP_PASSWORD", ""),
		SMTPFrom:     getEnv("SMTP_FROM", ""),

		ClaudeAPIKey:  getEnv("CLAUDE_API_KEY", ""),
		ClaudeModel:   getEnv("CLAUDE_MODEL", "claude-haiku-4-5"),

		GeminiAPIKey:          getEnv("GEMINI_API_KEY", ""),
		GeminiChatModel:       getEnv("GEMINI_CHAT_MODEL", "gemini-2.5-flash"),
		GeminiEmbedModel:      getEnv("GEMINI_EMBED_MODEL", "gemini-embedding-001"),
		AIChatHistoryMaxTurns: getEnvInt("AI_CHAT_MAX_TURNS", 12),
		AIChatHistoryMaxChars: getEnvInt("AI_CHAT_MAX_CHARS", 16000),
		AIRequestTimeoutSec:   getEnvInt("AI_REQUEST_TIMEOUT_SEC", 25),
		AIBackfillIntervalSec: getEnvInt("AI_BACKFILL_INTERVAL_SEC", 300),
		AIBackfillBatchSize:   getEnvInt("AI_BACKFILL_BATCH_SIZE", 25),

		FCMProjectID:          getEnv("FCM_PROJECT_ID", ""),
		FCMServiceAccountJSON: getEnv("FCM_SERVICE_ACCOUNT_JSON", ""),
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

func getEnvInt(key string, defaultValue int) int {
	if raw := os.Getenv(key); raw != "" {
		if v, err := strconv.Atoi(raw); err == nil {
			return v
		}
	}
	return defaultValue
}
