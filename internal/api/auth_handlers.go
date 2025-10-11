package api

import (
	"encoding/json"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"unicode"

	"github.com/arman300s/uni-portal/internal/models"
	"github.com/arman300s/uni-portal/pkg/auth"
	"github.com/arman300s/uni-portal/pkg/db"
	"gorm.io/gorm"
)

const (
	minPasswordLength = 8
	maxPasswordLength = 128
	maxNameLength     = 100
	maxEmailLength    = 255
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

type validationError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

// Error implements the error interface
func (e *validationError) Error() string {
	return e.Message
}

type errorResponse struct {
	Error   string            `json:"error"`
	Details []validationError `json:"details,omitempty"`
}

// validateEmail checks if email format is valid
func validateEmail(email string) error {
	if len(email) == 0 {
		return &validationError{Field: "email", Message: "email is required"}
	}
	if len(email) > maxEmailLength {
		return &validationError{Field: "email", Message: "email is too long"}
	}

	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)
	if !emailRegex.MatchString(email) {
		return &validationError{Field: "email", Message: "invalid email format"}
	}
	return nil
}

// validatePassword checks password strength
func validatePassword(password string) error {
	if len(password) == 0 {
		return &validationError{Field: "password", Message: "password is required"}
	}
	if len(password) < minPasswordLength {
		return &validationError{Field: "password", Message: "password must be at least 8 characters"}
	}
	if len(password) > maxPasswordLength {
		return &validationError{Field: "password", Message: "password is too long"}
	}

	var (
		hasUpper   bool
		hasLower   bool
		hasNumber  bool
		hasSpecial bool
	)

	for _, char := range password {
		switch {
		case unicode.IsUpper(char):
			hasUpper = true
		case unicode.IsLower(char):
			hasLower = true
		case unicode.IsNumber(char):
			hasNumber = true
		case unicode.IsPunct(char) || unicode.IsSymbol(char):
			hasSpecial = true
		}
	}

	if !hasUpper {
		return &validationError{Field: "password", Message: "password must contain at least one uppercase letter"}
	}
	if !hasLower {
		return &validationError{Field: "password", Message: "password must contain at least one lowercase letter"}
	}
	if !hasNumber {
		return &validationError{Field: "password", Message: "password must contain at least one number"}
	}
	if !hasSpecial {
		return &validationError{Field: "password", Message: "password must contain at least one special character"}
	}

	return nil
}

func validateName(name string) error {
	if len(name) == 0 {
		return &validationError{Field: "name", Message: "name is required"}
	}

	name = strings.TrimSpace(name)
	if len(name) == 0 {
		return &validationError{Field: "name", Message: "name cannot be empty"}
	}
	if len(name) > maxNameLength {
		return &validationError{Field: "name", Message: "name is too long"}
	}

	nameRegex := regexp.MustCompile(`^[a-zA-Z\s\-']+$`)
	if !nameRegex.MatchString(name) {
		return &validationError{Field: "name", Message: "name contains invalid characters"}
	}

	return nil
}

// sendJSONError sends a JSON error response
func sendJSONError(w http.ResponseWriter, message string, statusCode int, details []validationError) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	response := errorResponse{
		Error:   message,
		Details: details,
	}
	json.NewEncoder(w).Encode(response)
}

// SignupHandler @Summary User signup
// @Description Create a new user account
// @Tags auth
// @Accept json
// @Produce json
// @Param user body signupReq true "User signup data"
// @Success 201 {object} map[string]string
// @Failure 400 {object} errorResponse
// @Router /signup [post]
func SignupHandler(w http.ResponseWriter, r *http.Request) {
	var req signupReq

	r.Body = http.MaxBytesReader(w, r.Body, 1048576) // 1MB limit

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		sendJSONError(w, "Invalid request body", http.StatusBadRequest, nil)
		return
	}

	req.Email = strings.TrimSpace(strings.ToLower(req.Email))
	req.Name = strings.TrimSpace(req.Name)

	var validationErrors []validationError

	if err := validateName(req.Name); err != nil {
		if ve, ok := err.(*validationError); ok {
			validationErrors = append(validationErrors, *ve)
		}
	}

	if err := validateEmail(req.Email); err != nil {
		if ve, ok := err.(*validationError); ok {
			validationErrors = append(validationErrors, *ve)
		}
	}

	if err := validatePassword(req.Password); err != nil {
		if ve, ok := err.(*validationError); ok {
			validationErrors = append(validationErrors, *ve)
		}
	}

	if len(validationErrors) > 0 {
		sendJSONError(w, "Validation failed", http.StatusBadRequest, validationErrors)
		return
	}

	var userExists models.User
	if err := db.DB.Where("email = ?", req.Email).First(&userExists).Error; err == nil {
		sendJSONError(w, "Email already in use", http.StatusConflict, nil)
		return
	} else if err != gorm.ErrRecordNotFound {
		sendJSONError(w, "Database error", http.StatusInternalServerError, nil)
		return
	}

	hash, err := auth.HashPassword(req.Password)
	if err != nil {
		sendJSONError(w, "Failed to process password", http.StatusInternalServerError, nil)
		return
	}

	user := models.User{
		Name:     req.Name,
		Email:    req.Email,
		Password: hash,
	}

	if err := db.DB.Create(&user).Error; err != nil {
		sendJSONError(w, "Failed to create user", http.StatusInternalServerError, nil)
		return
	}

	token, err := auth.GenerateToken(user.ID)
	if err != nil {
		sendJSONError(w, "Failed to generate token", http.StatusInternalServerError, nil)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
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
// @Failure 401 {object} errorResponse
// @Router /login [post]
func LoginHandler(w http.ResponseWriter, r *http.Request) {
	var req loginReq

	r.Body = http.MaxBytesReader(w, r.Body, 1048576) // 1MB limit

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		sendJSONError(w, "Invalid request body", http.StatusBadRequest, nil)
		return
	}

	req.Email = strings.TrimSpace(strings.ToLower(req.Email))

	if req.Email == "" || req.Password == "" {
		sendJSONError(w, "Email and password are required", http.StatusBadRequest, nil)
		return
	}

	if err := validateEmail(req.Email); err != nil {
		sendJSONError(w, "Invalid credentials", http.StatusUnauthorized, nil)
		return
	}

	if len(req.Password) > maxPasswordLength {
		sendJSONError(w, "Invalid credentials", http.StatusUnauthorized, nil)
		return
	}

	var user models.User
	if err := db.DB.Where("email = ?", req.Email).First(&user).Error; err != nil {
		// Don't reveal whether user exists
		sendJSONError(w, "Invalid credentials", http.StatusUnauthorized, nil)
		return
	}

	if err := auth.CheckPassword(user.Password, req.Password); err != nil {
		sendJSONError(w, "Invalid credentials", http.StatusUnauthorized, nil)
		return
	}

	token, err := auth.GenerateToken(user.ID)
	if err != nil {
		sendJSONError(w, "Failed to generate token", http.StatusInternalServerError, nil)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"token": token,
		"id":    strconv.FormatUint(uint64(user.ID), 10),
	})
}
