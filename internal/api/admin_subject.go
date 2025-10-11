package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/arman300s/uni-portal/internal/models"
	"github.com/arman300s/uni-portal/pkg/db"
	"github.com/gorilla/mux"
	"gorm.io/gorm"
)

type SubjectRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	TeacherIDs  []uint `json:"teacher_ids"`
}

func AdminCreateSubjectHandler(w http.ResponseWriter, r *http.Request) {
	var req SubjectRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid JSON: "+err.Error(), http.StatusBadRequest)
		return
	}

	if req.Name == "" {
		http.Error(w, "subject name required", http.StatusBadRequest)
		return
	}

	subject := models.Subject{
		Name:        req.Name,
		Description: req.Description,
	}

	if len(req.TeacherIDs) > 0 {
		var teachers []models.User
		if err := db.DB.Where("id IN ?", req.TeacherIDs).Find(&teachers).Error; err != nil {
			http.Error(w, "failed to find teachers: "+err.Error(), http.StatusInternalServerError)
			return
		}

		for _, t := range teachers {
			if t.RoleID == nil || *t.RoleID != 2 {
				http.Error(w, fmt.Sprintf("user %d is not a teacher", t.ID), http.StatusBadRequest)
				return
			}
		}

		subject.Teachers = teachers
	}

	if err := db.DB.Create(&subject).Error; err != nil {
		http.Error(w, "failed to create subject: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "subject created successfully",
		"subject": subject,
	})
}

func AdminListSubjectsHandler(w http.ResponseWriter, r *http.Request) {
	var subjects []models.Subject
	if err := db.DB.Preload("Teachers").Find(&subjects).Error; err != nil {
		http.Error(w, "failed to fetch subjects: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(subjects)
}

func AdminGetSubjectHandler(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	var subject models.Subject
	if err := db.DB.Preload("Teachers").First(&subject, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			http.Error(w, "subject not found", http.StatusNotFound)
		} else {
			http.Error(w, "database error: "+err.Error(), http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(subject)
}

func AdminUpdateSubjectHandler(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]

	var req SubjectRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid JSON: "+err.Error(), http.StatusBadRequest)
		return
	}

	var subject models.Subject
	if err := db.DB.Preload("Teachers").First(&subject, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			http.Error(w, "subject not found", http.StatusNotFound)
		} else {
			http.Error(w, "database error: "+err.Error(), http.StatusInternalServerError)
		}
		return
	}

	if req.Name != "" {
		subject.Name = req.Name
	}
	if req.Description != "" {
		subject.Description = req.Description
	}

	if len(req.TeacherIDs) > 0 {
		var teachers []models.User
		if err := db.DB.Where("id IN ?", req.TeacherIDs).Find(&teachers).Error; err != nil {
			http.Error(w, "failed to fetch teachers: "+err.Error(), http.StatusInternalServerError)
			return
		}

		for _, t := range teachers {
			if t.RoleID == nil || *t.RoleID != 2 {
				http.Error(w, fmt.Sprintf("user %d is not a teacher", t.ID), http.StatusBadRequest)
				return
			}
		}

		if err := db.DB.Model(&subject).Association("Teachers").Replace(teachers); err != nil {
			http.Error(w, "failed to update teachers: "+err.Error(), http.StatusInternalServerError)
			return
		}
	}

	if err := db.DB.Save(&subject).Error; err != nil {
		http.Error(w, "failed to update subject: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "subject updated successfully",
	})
}

func AdminDeleteSubjectHandler(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	sid, err := strconv.Atoi(id)
	if err != nil {
		http.Error(w, "invalid subject id", http.StatusBadRequest)
		return
	}

	if err := db.DB.Delete(&models.Subject{}, sid).Error; err != nil {
		http.Error(w, "failed to delete subject: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"message": "subject deleted successfully",
	})
}
