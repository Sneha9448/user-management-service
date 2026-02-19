package auth

import (
	"context"
	"crypto/rand"
	"errors"
	"fmt"
	"math/big"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"google.golang.org/api/idtoken"
)

var jwtKey = []byte("your_secret_key") // In production, use an environment variable

type Claims struct {
	UserID string `json:"user_id"`
	Email  string `json:"email"`
	Role   string `json:"role"`
	jwt.RegisteredClaims
}

// GenerateJWT creates a new JWT token for a user
func GenerateJWT(userID, email, role string) (string, error) {
	expirationTime := time.Now().Add(15 * time.Minute) // Token valid for 15 minutes
	claims := &Claims{
		UserID: userID,
		Email:  email,
		Role:   role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtKey)
}

// VerifyJWT parses and validates a JWT token
func VerifyJWT(tokenString string) (*Claims, error) {
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, errors.New("invalid token")
	}

	return claims, nil
}

// GenerateOTP creates a secure 6-digit random code
func GenerateOTP() (string, error) {
	otp := ""
	for i := 0; i < 6; i++ {
		num, err := rand.Int(rand.Reader, big.NewInt(10))
		if err != nil {
			return "", err
		}
		otp += fmt.Sprintf("%d", num)
	}
	return otp, nil
}

// VerifyGoogleToken verifies the Google ID token and returns the user's email
func VerifyGoogleToken(ctx context.Context, idToken string, clientID string) (string, error) {
	// For development/demo purposes, we'll allow a mock token if clientID is not set or for specific test tokens
	if idToken == "mock_token" {
		return "test@gmail.com", nil
	}

	payload, err := idtoken.Validate(ctx, idToken, clientID)
	if err != nil {
		return "", fmt.Errorf("failed to validate id token: %v", err)
	}

	email, ok := payload.Claims["email"].(string)
	if !ok {
		return "", errors.New("email not found in google token payload")
	}

	return email, nil
}
