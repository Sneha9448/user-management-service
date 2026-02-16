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

// SaveOTP stores a new OTP in the database
func SaveOTP(otp *models.OTP) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if database.DB == nil {
		return errors.New("database connection is not initialized")
	}

	query := `INSERT INTO otps (email, otp, expires_at) VALUES ($1, $2, $3) RETURNING id`

	err := database.DB.QueryRow(ctx, query, otp.Email, otp.OTP, otp.ExpiresAt).Scan(&otp.ID)
	if err != nil {
		log.Printf("Error saving OTP: %v", err)
		return err
	}
	return nil
}

// GetLatestOTP retrieves the most recent OTP for an email
func GetLatestOTP(email string) (*models.OTP, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if database.DB == nil {
		return nil, errors.New("database connection is not initialized")
	}

	query := `SELECT id, email, otp, expires_at, is_used, attempt_count, created_at FROM otps 
			  WHERE email = $1 ORDER BY created_at DESC LIMIT 1`

	var otp models.OTP
	err := database.DB.QueryRow(ctx, query, email).Scan(
		&otp.ID, &otp.Email, &otp.OTP, &otp.ExpiresAt, &otp.IsUsed, &otp.AttemptCount, &otp.CreatedAt,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil // No OTP found
		}
		log.Printf("Error fetching OTP: %v", err)
		return nil, err
	}
	return &otp, nil
}

// IncrementOTPAttempts increases the attempt count for an OTP
func IncrementOTPAttempts(id int) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if database.DB == nil {
		return errors.New("database connection is not initialized")
	}

	query := `UPDATE otps SET attempt_count = attempt_count + 1 WHERE id = $1`
	_, err := database.DB.Exec(ctx, query, id)
	return err
}

// MarkOTPAsUsed marks an OTP as used
func MarkOTPAsUsed(id int) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if database.DB == nil {
		return errors.New("database connection is not initialized")
	}

	query := `UPDATE otps SET is_used = TRUE WHERE id = $1`
	_, err := database.DB.Exec(ctx, query, id)
	return err
}
