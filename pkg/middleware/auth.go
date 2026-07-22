package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/s1ntezc0der/bazis-restapi/pkg/errors"
	"github.com/s1ntezc0der/bazis-restapi/pkg/jwt"
)

type contextKey string

const UserIDKey contextKey = "user_id"

func AuthMiddleware(jwtService *jwt.JWTConfig) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				http.Error(w, errors.ErrUnauthorized.Error(), http.StatusUnauthorized)
				return
			}

			parts := strings.Split(authHeader, " ")
			if len(parts) != 2 || parts[0] != "Bearer" {
				http.Error(w, errors.ErrInvalidToken.Error(), http.StatusUnauthorized)
				return
			}

			claims, err := jwtService.Validate(parts[1])
			if err != nil {
				http.Error(w, errors.ErrInvalidToken.Error(), http.StatusUnauthorized)
				return
			}

			ctx := context.WithValue(r.Context(), UserIDKey, claims.UserID)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

