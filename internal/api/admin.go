package api

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"strings"

	"github.com/arman300s/uni-portal/internal/models"
	"github.com/arman300s/uni-portal/pkg/db"
	"github.com/gorilla/mux"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type CreateUserRequest struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
	RoleName string `json:"role"`
}

type UpdateUserRequest struct {
	Email    string `json:"email"`
	RoleName string `json:"role"`
}

func AdminListUsersHandler(w http.ResponseWriter, r *http.Request) {
	var users []models.User
	if err := db.DB.Preload("Role").Find(&users).Error; err != nil {
		http.Error(w, "failed to fetch users: "+err.Error(), http.StatusInternalServerError)
		return
	}

	var result []map[string]interface{}
	for _, u := range users {
		result = append(result, map[string]interface{}{
			"id":    u.ID,
			"name":  u.Name,
			"email": u.Email,
			"role":  u.Role.Name,
		})
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}

func AdminGetUserHandler(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	var user models.User
	if err := db.DB.Preload("Role").First(&user, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			http.Error(w, "user not found", http.StatusNotFound)
		} else {
			http.Error(w, "failed to fetch user: "+err.Error(), http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"id":    user.ID,
		"name":  user.Name,
		"email": user.Email,
		"role":  user.Role.Name,
	})
}

func AdminUpdateUserHandler(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]

	var req UpdateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid JSON: "+err.Error(), http.StatusBadRequest)
		return
	}

	var user models.User
	if err := db.DB.First(&user, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			http.Error(w, "user not found", http.StatusNotFound)
		} else {
			http.Error(w, "database error: "+err.Error(), http.StatusInternalServerError)
		}
		return
	}

	if req.RoleName != "" {
		req.RoleName = strings.ToLower(req.RoleName)
		var role models.Role
		if err := db.DB.First(&role, "name = ?", req.RoleName).Error; err != nil {
			http.Error(w, "invalid role: "+req.RoleName, http.StatusBadRequest)
			return
		}
		user.RoleID = &role.ID
	}

	if err := db.DB.Save(&user).Error; err != nil {
		http.Error(w, "failed to update user: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "user updated successfully",
	})
}

func AdminCreateUserHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req CreateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid JSON: "+err.Error(), http.StatusBadRequest)
		return
	}

	req.RoleName = strings.ToLower(req.RoleName)

	var role models.Role
	if err := db.DB.First(&role, "name = ?", req.RoleName).Error; err != nil {
		http.Error(w, "invalid role: "+req.RoleName, http.StatusBadRequest)
		return
	}

	var existingUser models.User
	if err := db.DB.First(&existingUser, "email = ?", req.Email).Error; err == nil {
		http.Error(w, "email already in use", http.StatusConflict)
		return
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		http.Error(w, "database error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, "failed to hash password", http.StatusInternalServerError)
		return
	}

	user := models.User{
		Name:     req.Name,
		Email:    req.Email,
		Password: string(hashedPassword),
		RoleID:   &role.ID,
	}

	if err := db.DB.Create(&user).Error; err != nil {
		http.Error(w, "failed to create user: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "user created successfully",
		"user": map[string]interface{}{
			"id":    user.ID,
			"name":  user.Name,
			"email": user.Email,
			"role":  role.Name,
		},
	})
}

func AdminDeleteUserHandler(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	uid, err := strconv.Atoi(id)
	if err != nil {
		http.Error(w, "invalid user id", http.StatusBadRequest)
		return
	}

	var user models.User
	if err := db.DB.First(&user, uid).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			http.Error(w, "user not found", http.StatusNotFound)
			return
		}
		http.Error(w, "failed to check user: "+err.Error(), http.StatusInternalServerError)
		return
	}

	if err := db.DB.Delete(&user).Error; err != nil {
		http.Error(w, "failed to delete user: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"message": "user deleted successfully",
	})
}
