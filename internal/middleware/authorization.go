package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/faqq11/lib-management/internal/helper"
)

type contextKey string

const UserContextKey contextKey = "user"

type UserClaims struct {
	UserID   int    `json:"user_id"`
	Username string `json:"username"`
	Role     string `json:"role"`
}

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		authHeader := request.Header.Get("Authorization")
		if authHeader == "" {
			helper.ErrorResponse(writer, http.StatusUnauthorized, "Authorization header required")
			return
		}

		if !strings.HasPrefix(authHeader, "Bearer ") {
			helper.ErrorResponse(writer, http.StatusUnauthorized, "Invalid authorization format")
			return
		}

		token := strings.TrimPrefix(authHeader, "Bearer ")
		if token == "" {
			helper.ErrorResponse(writer, http.StatusUnauthorized, "Token required")
			return
		}

		claims, err := helper.VerifyJWT(token)
		if err != nil {
			helper.ErrorResponse(writer, http.StatusUnauthorized, "Invalid token")
			return
		}

		userID, ok := claims["userId"].(float64)
		if !ok {
			helper.ErrorResponse(writer, http.StatusUnauthorized, "Invalid token format")
			return
		}

		username, ok := claims["username"].(string)
		if !ok {
			helper.ErrorResponse(writer, http.StatusUnauthorized, "Invalid token format")
			return
		}

		role, ok := claims["role"].(string)
		if !ok {
			helper.ErrorResponse(writer, http.StatusUnauthorized, "Invalid token format")
			return
		}

		userClaims := UserClaims{
			UserID:   int(userID),
			Username: username,
			Role:     role,
		}

		ctx := context.WithValue(request.Context(), UserContextKey, userClaims)
		next.ServeHTTP(writer, request.WithContext(ctx))
	})
}

func AdminOnlyMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		user := request.Context().Value(UserContextKey)
		if user == nil {
			helper.ErrorResponse(writer, http.StatusUnauthorized, "User context not found")
			return
		}

		userClaims := user.(UserClaims)
		if userClaims.Role != "admin" {
			helper.ErrorResponse(writer, http.StatusForbidden, "Admin access required")
			return
		}

		next.ServeHTTP(writer, request)
	})
}