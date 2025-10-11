package api

import (
	"encoding/json"
	"net/http"

	"github.com/arman300s/uni-portal/internal/models"
	"github.com/arman300s/uni-portal/pkg/db"
	"github.com/arman300s/uni-portal/pkg/middleware"
)

func TeacherListMySubjectsHandler(w http.ResponseWriter, r *http.Request) {
	uid, ok := middleware.UserIDFromContext(r.Context())
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	var teacher models.User
	if err := db.DB.First(&teacher, uid).Error; err != nil {
		http.Error(w, "user not found", http.StatusNotFound)
		return
	}

	var subjects []models.Subject
	if err := db.DB.Preload("Teachers").
		Joins("JOIN subject_teachers st ON st.subject_id = subjects.id").
		Where("st.user_id = ?", uid).
		Find(&subjects).Error; err != nil {
		http.Error(w, "failed to fetch subjects: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(subjects)
}
