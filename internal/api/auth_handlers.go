package api

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/arman300s/uni-portal/internal/models"
	"github.com/arman300s/uni-portal/pkg/auth"
	"github.com/arman300s/uni-portal/pkg/db"
)

type signupReq struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type loginReq struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// SignupHandler @Summary User signup
// @Description Create a new user account
// @Tags auth
// @Accept json
// @Produce json
// @Param user body signupReq true "User signup data"
// @Success 201 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Router /signup [post]
func SignupHandler(w http.ResponseWriter, r *http.Request) {
	var req signupReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid body", http.StatusBadRequest)
		return
	}

	if req.Email == "" || req.Password == "" {
		http.Error(w, "Email and password required", http.StatusBadRequest)
		return
	}

	var existingUser models.User
	if err := db.DB.Where("email = ?", req.Email).First(&existingUser).Error; err != nil {
		http.Error(w, "Email already in use", http.StatusNotFound)
	}
	hash, err := auth.HashPassword(req.Password)
	if err != nil {
		http.Error(w, "failed to hash", http.StatusBadRequest)
		return
	}
	user := models.User{
		Name:     req.Name,
		Email:    req.Email,
		Password: hash,
	}
	if err := db.DB.Create(&user).Error; err != nil {
		http.Error(w, "failed to create user", http.StatusBadRequest)
		return
	}
	token, err := auth.GenerateToken(user.ID)
	if err != nil {
		http.Error(w, "failed to generate token", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"token": token,
		"id":    strconv.FormatUint(uint64(user.ID), 10),
	})
}

// LoginHandler @Summary User login
// @Description Authenticate user and return JWT token
// @Tags auth
// @Accept json
// @Produce json
// @Param credentials body loginReq true "User credentials"
// @Success 200 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Router /login [post]
func LoginHandler(w http.ResponseWriter, r *http.Request) {
	var req loginReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid body", http.StatusBadRequest)
		return
	}
	var user models.User
	if err := db.DB.Where("email = ?", req.Email).First(&user).Error; err != nil {
		http.Error(w, "invalid credentials", http.StatusUnauthorized)
		return
	}
	if err := auth.CheckPassword(user.Password, req.Password); err != nil {
		http.Error(w, "invalid credentials", http.StatusUnauthorized)
		return
	}
	token, err := auth.GenerateToken(user.ID)
	if err != nil {
		http.Error(w, "failed to generate token", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"token": token,
	})
}
