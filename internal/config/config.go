package config

import (
	"os"
)

type Config struct {
	DatabaseURL string
	Port        string
}

func LoadConfig() *Config {
	return &Config{
		DatabaseURL: getEnv(
			"DATABASE_URL",
			"postgres://postgres:postgres@localhost:5432/user?sslmode=disable&search_path=public",
		),
		Port: getEnv("PORT", "8081"),
	}
}

func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}
