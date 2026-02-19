package main

import (
	"context"
	"fmt"
	"log"
	"user-management-service/internal/config"
	"user-management-service/internal/database"
)

func main() {
	// Load config
	cfg := config.LoadConfig()
	if cfg.DatabaseURL == "" {
		cfg.DatabaseURL = "postgres://postgres:postgres@localhost:5432/user?sslmode=disable"
	}

	fmt.Printf("Connecting to %s...\n", cfg.DatabaseURL)
	database.ConnectDB(cfg.DatabaseURL)
	defer database.CloseDB()

	ctx := context.Background()

	// 1. List all users
	fmt.Println("\n--- Current Users in Database ---")
	rows, err := database.DB.Query(ctx, "SELECT id, name, email, role FROM users")
	if err != nil {
		log.Fatalf("Failed to query users: %v", err)
	}
	defer rows.Close()

	var usersExist bool
	for rows.Next() {
		usersExist = true
		var id int
		var name, email, role string
		err := rows.Scan(&id, &name, &email, &role)
		if err != nil {
			log.Printf("Error scanning row: %v", err)
			continue
		}
		fmt.Printf("ID: %d | Name: %s | Email: %s | Role: %s\n", id, name, email, role)
	}

	if !usersExist {
		fmt.Println("No users found in database.")
	}

	// 2. Automatically promote every user to ADMIN for testing if requested,
	// or provide instructions.
	fmt.Println("\n--- Updating all users to ADMIN for debugging ---")
	tag, err := database.DB.Exec(ctx, "UPDATE users SET role = 'ADMIN'")
	if err != nil {
		log.Fatalf("Failed to update users: %v", err)
	}
	fmt.Printf("Updated %d users to ADMIN.\n", tag.RowsAffected())

	fmt.Println("\nDone! Please RESTART your server and LOG OUT/IN of the web app.")
}
