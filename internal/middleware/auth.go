package middleware

import (
	"context"
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
				// We don't block here because some queries might be public
				// The resolver will check if the user is authenticated if needed
				next.ServeHTTP(w, r)
				return
			}

			// Put user info into context
			user := &User{
				ID:    claims.UserID,
				Email: claims.Email,
			}
			ctx := context.WithValue(r.Context(), UserCtxKey, user)

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
