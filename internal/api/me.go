package api

import (
	"encoding/json"
	"net/http"

	"github.com/arman300s/uni-portal/internal/models"
	"github.com/arman300s/uni-portal/pkg/db"
	"github.com/arman300s/uni-portal/pkg/middleware"
)

// MeHandler @Summary Get current user
// @Description Get current authenticated user information
// @Tags auth
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {object} models.User
// @Failure 401 {object} map[string]string
// @Router /me [get]
func MeHandler(w http.ResponseWriter, r *http.Request) {
	uid, ok := middleware.UserIDFromContext(r.Context())
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}
	var user models.User
	if err := db.DB.First(&user, uid).Error; err != nil {
		http.Error(w, "user not found", http.StatusNotFound)
		return
	}
	user.Password = ""
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}
