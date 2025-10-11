package api

import (
	"encoding/json"
	"net/http"

	"github.com/arman300s/uni-portal/internal/models"
	"github.com/arman300s/uni-portal/pkg/db"
)

func StudentListSubjectsHandler(w http.ResponseWriter, r *http.Request) {
	var subjects []models.Subject
	if err := db.DB.Preload("Teachers").Find(&subjects).Error; err != nil {
		http.Error(w, "failed to fetch subjects: "+err.Error(), http.StatusInternalServerError)
		return
	}

	type SubjectResponse struct {
		ID          uint     `json:"id"`
		Name        string   `json:"name"`
		Description string   `json:"description"`
		Teachers    []string `json:"teachers"`
	}

	var resp []SubjectResponse
	for _, s := range subjects {
		var teacherNames []string
		for _, t := range s.Teachers {
			teacherNames = append(teacherNames, t.Name)
		}
		resp = append(resp, SubjectResponse{
			ID:          s.ID,
			Name:        s.Name,
			Description: s.Description,
			Teachers:    teacherNames,
		})
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}
