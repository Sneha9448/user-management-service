package repository

import (
	"context"
	"errors"
	"log"
	"time"

	"user-management-service/internal/database"
	"user-management-service/internal/models"

	"github.com/jackc/pgx/v5"
)

// CreateUser inserts a new user into the database
func CreateUser(user *models.User) error {
	// Use a 5-second timeout for the database operation
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if database.DB == nil {
		return errors.New("database connection is not initialized")
	}

	query := `INSERT INTO users (name, email, role) VALUES ($1, $2, $3) RETURNING id`

	err := database.DB.QueryRow(ctx, query, user.Name, user.Email, user.Role).Scan(&user.ID)
	if err != nil {
		log.Printf("Error creating user: %v", err)
		return err
	}
	return nil
}

// GetUserByID fetches a user by their ID
func GetUserByID(id int) (*models.User, error) {
	// Use a 5-second timeout for the database operation
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if database.DB == nil {
		return nil, errors.New("database connection is not initialized")
	}

	query := `SELECT id, name, email, role FROM users WHERE id = $1`

	var user models.User
	err := database.DB.QueryRow(ctx, query, id).Scan(&user.ID, &user.Name, &user.Email, &user.Role)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil // User not found
		}
		log.Printf("Error fetching user: %v", err)
		return nil, err
	}
	return &user, nil
}

// GetAllUsers retrieves all users from the database
func GetAllUsers() ([]*models.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if database.DB == nil {
		return nil, errors.New("database connection is not initialized")
	}

	query := `SELECT id, name, email, role FROM users`

	rows, err := database.DB.Query(ctx, query)
	if err != nil {
		log.Printf("Error querying users: %v", err)
		return nil, err
	}
	defer rows.Close()

	var users []*models.User
	for rows.Next() {
		var user models.User
		if err := rows.Scan(&user.ID, &user.Name, &user.Email, &user.Role); err != nil {
			log.Printf("Error scanning user row: %v", err)
			return nil, err
		}
		users = append(users, &user)
	}

	if err := rows.Err(); err != nil {
		log.Printf("Error iterating user rows: %v", err)
		return nil, err
	}

	return users, nil
}

// UpdateUser updates an existing user's information
func UpdateUser(user *models.User) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if database.DB == nil {
		return errors.New("database connection is not initialized")
	}

	query := `UPDATE users SET name = $1, email = $2, role = $3 WHERE id = $4 returning id`

	err := database.DB.QueryRow(ctx, query, user.Name, user.Email, user.Role, user.ID).Scan(&user.ID)
	if err != nil {
		if err == pgx.ErrNoRows {
			return errors.New("user not found")
		}
		log.Printf("Error updating user: %v", err)
		return err
	}
	return nil
}

// DeleteUser removes a user from the database
func DeleteUser(id int) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if database.DB == nil {
		return errors.New("database connection is not initialized")
	}

	query := `DELETE FROM users WHERE id = $1`

	result, err := database.DB.Exec(ctx, query, id)
	if err != nil {
		log.Printf("Error deleting user: %v", err)
		return err
	}

	if result.RowsAffected() == 0 {
		return errors.New("user not found")
	}

	return nil
}

// GetUserByEmail fetches a user by their email
func GetUserByEmail(email string) (*models.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if database.DB == nil {
		return nil, errors.New("database connection is not initialized")
	}

	query := `SELECT id, name, email, role FROM users WHERE email = $1`

	var user models.User
	err := database.DB.QueryRow(ctx, query, email).Scan(&user.ID, &user.Name, &user.Email, &user.Role)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, errors.New("user not found")
		}
		log.Printf("Error fetching user by email: %v", err)
		return nil, err
	}
	return &user, nil
}
