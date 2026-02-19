package middleware

import (
	"context"
	"log"
	"net/http"
	"strings"
	"user-management-service/internal/auth"
)

type contextKey struct {
	name string
}

var UserCtxKey = &contextKey{"user"}

// User struct to be stored in context
type User struct {
	ID    string
	Email string
	Role  string
}

// AuthMiddleware extracts the user from the JWT in the Authorization header
func AuthMiddleware() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			header := r.Header.Get("Authorization")

			// Allow unauthenticated users in (resolvers will check for user later)
			if header == "" {
				next.ServeHTTP(w, r)
				return
			}

			// Validate the token
			tokenStr := strings.TrimPrefix(header, "Bearer ")
			claims, err := auth.VerifyJWT(tokenStr)
			if err != nil {
				log.Printf("Auth Error: JWT verification failed: %v", err)
				next.ServeHTTP(w, r)
				return
			}

			// Put user info into context
			user := &User{
				ID:    claims.UserID,
				Email: claims.Email,
				Role:  claims.Role,
			}

			ctx := context.WithValue(r.Context(), UserCtxKey, user)

			log.Printf("Auth Success: User %s with Role %s identified", user.Email, user.Role)
			// Continue with the new context
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// ForContext finds the user from the context. REQUIRES Middleware to have run.
func ForContext(ctx context.Context) *User {
	raw, _ := ctx.Value(UserCtxKey).(*User)
	return raw
}
