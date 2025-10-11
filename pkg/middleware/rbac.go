package middleware

import (
	"context"
	"net/http"

	"github.com/arman300s/uni-portal/internal/models"
	"github.com/arman300s/uni-portal/pkg/db"
)

type ctxUserKey string

const userCtxKey ctxUserKey = "user"

func LoadUserMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userID, ok := UserIDFromContext(r.Context())
		if !ok {
			http.Error(w, "user id not found in context", http.StatusUnauthorized)
			return
		}

		var user models.User
		if err := db.DB.Preload("Role").First(&user, userID).Error; err != nil {
			http.Error(w, "user not found", http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), userCtxKey, user)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func RequireRole(allowedRoles ...string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			userVal := r.Context().Value(userCtxKey)
			if userVal == nil {
				http.Error(w, "user not in context", http.StatusUnauthorized)
				return
			}
			user := userVal.(models.User)

			for _, role := range allowedRoles {
				if user.Role.Name == role {
					next.ServeHTTP(w, r)
					return
				}
			}

			http.Error(w, "forbidden: insufficient role", http.StatusForbidden)
		})
	}
}
