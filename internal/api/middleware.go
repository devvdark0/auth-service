package api

import (
	"context"
	"log/slog"
	"net/http"
	"strings"

	"github.com/devvdark0/auth-service/internal/auth"
)

func LoggingMiddleware(log *slog.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			log.Debug("Request URI: %s, Method: %s", r.RequestURI, r.Method)
			next.ServeHTTP(w, r)
		})
	}
}

func AuthMiddleware(jwt *auth.JWTManager) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				http.Error(w, "Authorization header required", http.StatusUnauthorized)
				return
			}

			parts := strings.Split(authHeader, " ")

			if len(parts) != 2 && parts[0] != "Bearer" {
				http.Error(w, "Invalid authorization foramt", http.StatusUnauthorized)
				return
			}

			token := parts[1]

			claims, err := jwt.ValidateToken(token)
			if err != nil {
				http.Error(w, "Invalid or expired token", http.StatusUnauthorized)
				return
			}

			userId, ok := claims["sub"].(string)
			if !ok {
				http.Error(w, "Invalid token claims", http.StatusUnauthorized)
				return
			}

			ctx := context.WithValue(r.Context(), "userId", userId)

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}

}
