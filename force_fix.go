package main

import (
	"context"
	"fmt"
	"os"
	"user-management-service/internal/config"
	"user-management-service/internal/database"
)

func main() {
	f, _ := os.Create("db_fix_log.txt")
	defer f.Close()

	fmt.Fprintf(f, "Starting DB fix...\n")
	cfg := config.LoadConfig()
	database.ConnectDB(cfg.DatabaseURL)

	ctx := context.Background()

	// Promote everyone to ADMIN for now
	tag, err := database.DB.Exec(ctx, "UPDATE users SET role = 'ADMIN'")
	if err != nil {
		fmt.Fprintf(f, "ERROR updating: %v\n", err)
		return
	}

	fmt.Fprintf(f, "SUCCESS: Updated %d users to ADMIN.\n", tag.RowsAffected())

	// List users
	rows, _ := database.DB.Query(ctx, "SELECT email, role FROM users")
	for rows.Next() {
		var e, r string
		rows.Scan(&e, &r)
		fmt.Fprintf(f, "User: %s | Role: %s\n", e, r)
	}
}
