package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"user-management-service/internal/config"
	"user-management-service/internal/database"
)

func main() {
	cfg := config.LoadConfig()
	database.ConnectDB(cfg.DatabaseURL)
	defer database.CloseDB()

	migration, err := os.ReadFile("migrations/20260212_create_otps.sql")
	if err != nil {
		log.Fatalf("Failed to read migration: %v", err)
	}

	_, err = database.DB.Exec(context.Background(), string(migration))
	if err != nil {
		log.Fatalf("Failed to run migration: %v", err)
	}

	fmt.Println("Migration applied successfully!")
}
