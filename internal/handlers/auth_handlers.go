package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"user-management-service/internal/auth"
	"user-management-service/internal/email"
	"user-management-service/internal/models"
	"user-management-service/internal/repository"
)

// RequestOTP handles the request to generate and send an OTP
func RequestOTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	defer r.Body.Close()

	var payload struct {
		Email string `json:"email"`
	}

	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		http.Error(w, `{"error": "Invalid request payload"}`, http.StatusBadRequest)
		return
	}

	if payload.Email == "" {
		http.Error(w, `{"error": "Email is required"}`, http.StatusBadRequest)
		return
	}

	// 1. Generate OTP
	otpCode, err := auth.GenerateOTP()
	if err != nil {
		log.Printf("Failed to generate OTP: %v", err)
		http.Error(w, `{"error": "Internal server error"}`, http.StatusInternalServerError)
		return
	}

	// 2. Save OTP to Database
	otp := &models.OTP{
		Email:     payload.Email,
		OTP:       otpCode,
		ExpiresAt: time.Now().Add(10 * time.Minute),
	}

	if err := repository.SaveOTP(otp); err != nil {
		log.Printf("Failed to save OTP: %v", err)
		http.Error(w, `{"error": "Failed to process request"}`, http.StatusInternalServerError)
		return
	}

	// 3. Send OTP Email
	if err := email.SendOTPEmail(payload.Email, otpCode); err != nil {
		log.Printf("Failed to send OTP email to %s: %v", payload.Email, err)
		// Note: OTP is saved but email failed. Client can retry.
		http.Error(w, `{"error": "Failed to send OTP email"}`, http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "OTP sent successfully"})
}

// VerifyOTP handles OTP validation and JWT generation
func VerifyOTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	defer r.Body.Close()

	var payload struct {
		Email string `json:"email"`
		OTP   string `json:"otp"`
	}

	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		http.Error(w, `{"error": "Invalid request payload"}`, http.StatusBadRequest)
		return
	}

	if payload.Email == "" || payload.OTP == "" {
		http.Error(w, `{"error": "Email and OTP are required"}`, http.StatusBadRequest)
		return
	}

	// 1. Get Latest OTP
	storedOTP, err := repository.GetLatestOTP(payload.Email)
	if err != nil {
		log.Printf("Database error fetching OTP: %v", err)
		http.Error(w, `{"error": "Internal server error"}`, http.StatusInternalServerError)
		return
	}

	if storedOTP == nil {
		http.Error(w, `{"error": "Invalid or expired OTP"}`, http.StatusUnauthorized)
		return
	}

	// 2. Validate OTP
	if storedOTP.IsUsed {
		http.Error(w, `{"error": "OTP already used"}`, http.StatusUnauthorized)
		return
	}

	if storedOTP.ExpiresAt.Before(time.Now()) {
		http.Error(w, `{"error": "OTP expired"}`, http.StatusUnauthorized)
		return
	}

	if storedOTP.OTP != payload.OTP {
		// Increment attempts
		repository.IncrementOTPAttempts(storedOTP.ID)
		http.Error(w, `{"error": "Invalid OTP"}`, http.StatusUnauthorized)
		return
	}

	// 3. Mark OTP as used
	if err := repository.MarkOTPAsUsed(storedOTP.ID); err != nil {
		log.Printf("Failed to mark OTP as used: %v", err)
	}

	// 4. Generate JWT Token
	role := models.RoleUser
	user, err := repository.GetUserByEmail(payload.Email)
	if err == nil && user != nil {
		role = user.Role
	}

	token, err := auth.GenerateJWT(payload.Email, payload.Email, role)
	if err != nil {
		log.Printf("Failed to generate JWT: %v", err)
		http.Error(w, `{"error": "Failed to generate token"}`, http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"token":   token,
		"message": "Login successful",
	})
}
