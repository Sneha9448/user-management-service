package repository

import (
	"fmt"
	"testing"
	"user-management-service/internal/database"
	"user-management-service/internal/models"
)

// BenchmarkCreateUser benchmarks the CreateUser function
func BenchmarkCreateUser(b *testing.B) {
	// Setup: Ensure DB is connected
	// In a real scenario, we'd use a test DB or mock, but here we want to measure real DB impact as per user request
	databaseURL := "postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable"
	database.ConnectDB(databaseURL)
	defer database.CloseDB()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		user := &models.User{
			Name:  fmt.Sprintf("Bench User %d-%d", b.N, i),
			Email: fmt.Sprintf("bench%d-%d@example.com", b.N, i),
		}
		err := CreateUser(user)
		if err != nil {
			b.Fatalf("Failed to create user: %v", err)
		}
	}
}

// BenchmarkGetAllUsers benchmarks retrieval
func BenchmarkGetAllUsers(b *testing.B) {
	databaseURL := "postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable"
	database.ConnectDB(databaseURL)
	defer database.CloseDB()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := GetAllUsers()
		if err != nil {
			b.Fatalf("Failed to get all users: %v", err)
		}
	}
}
