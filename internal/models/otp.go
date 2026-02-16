package models

import "time"

type OTP struct {
	ID           int       `json:"id"`
	Email        string    `json:"email"`
	OTP          string    `json:"otp"`
	ExpiresAt    time.Time `json:"expires_at"`
	IsUsed       bool      `json:"is_used"`
	AttemptCount int       `json:"attempt_count"`
	CreatedAt    time.Time `json:"created_at"`
}
