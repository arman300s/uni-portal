package middleware

import (
	"context"
	"github.com/arman300s/uni-portal/pkg/auth"
	"net/http"
	"strings"
)

type ctxKey string

const userIDKey ctxKey = "userID"

func UserIDFromContext(ctx context.Context) (uint, bool) {
	v := ctx.Value(userIDKey)
	if v == nil {
		return 0, false
	}
	id, ok := v.(uint)
	return id, ok
}

func JWTAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		h := r.Header.Get("Authorization")
		if h == "" {
			http.Error(w, "missing auth header", http.StatusUnauthorized)
			return
		}
		parts := strings.SplitN(h, " ", 2)
		if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
			http.Error(w, "invalid auth header", http.StatusUnauthorized)
			return
		}
		token := parts[1]
		claims, err := auth.ParseToken(token)
		if err != nil {
			http.Error(w, "invalid token: "+err.Error(), http.StatusUnauthorized)
			return
		}
		ctx := context.WithValue(r.Context(), userIDKey, claims.UserID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
