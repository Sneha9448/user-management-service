package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DatabaseURL  string
	Port         string
	SMTPHost     string
	SMTPPort     string
	SMTPEmail    string
	SMTPPassword string
	JWTSecret    string
}

func LoadConfig() *Config {
	// Load .env file (try CWD and parent directories)
	if err := godotenv.Load(); err != nil {
		if err := godotenv.Load("../.env"); err != nil {
			godotenv.Load("../../.env")
		}
	}

	return &Config{
		DatabaseURL: getEnv(
			"DATABASE_URL",
			"postgres://postgres:postgres@localhost:5432/user?sslmode=disable&search_path=public",
		),
		Port:         getEnv("PORT", "8080"),
		SMTPHost:     getEnv("SMTP_HOST", "smtp.gmail.com"),
		SMTPPort:     getEnv("SMTP_PORT", "587"),
		SMTPEmail:    getEnv("SMTP_EMAIL", ""),
		SMTPPassword: getEnv("SMTP_PASSWORD", ""),
		JWTSecret:    getEnv("JWT_SECRET", "your_secret_key"),
	}
}

func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}
