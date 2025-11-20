package controllers

import (
	"encoding/json"
	"net/http"

	"github.com/arman300s/uni-portal/internal/core/contracts"
	"github.com/arman300s/uni-portal/internal/core/services"
)

// AuthController handles signup/login requests.
type AuthController struct {
	service *services.AuthService
}

func NewAuthController(service *services.AuthService) *AuthController {
	return &AuthController{service: service}
}

// Signup godoc
// @Summary User signup
// @Description Create a new user account
// @Tags auth
// @Accept json
// @Produce json
// @Param signup body contracts.SignupInput true "Signup data"
// @Success 201 {object} contracts.AuthResponse
// @Failure 400 {object} ErrorResponse
// @Failure 409 {object} ErrorResponse
// @Router /signup [post]
func (c *AuthController) Signup(w http.ResponseWriter, r *http.Request) {
	var input contracts.SignupInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body", nil)
		return
	}

	resp, err := c.service.Signup(r.Context(), input)
	if err != nil {
		handleAuthError(w, err)
		return
	}

	writeJSON(w, http.StatusCreated, resp)
}

// Login godoc
// @Summary User login
// @Description Authenticate user and return JWT token
// @Tags auth
// @Accept json
// @Produce json
// @Param credentials body contracts.LoginInput true "User credentials"
// @Success 200 {object} contracts.AuthResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Router /login [post]
func (c *AuthController) Login(w http.ResponseWriter, r *http.Request) {
	var input contracts.LoginInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body", nil)
		return
	}

	resp, err := c.service.Login(r.Context(), input)
	if err != nil {
		handleAuthError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, resp)
}

func handleAuthError(w http.ResponseWriter, err error) {
	switch e := err.(type) {
	case contracts.ValidationErrors:
		writeError(w, http.StatusBadRequest, "validation failed", e)
	default:
		switch err {
		case contracts.ErrEmailInUse:
			writeError(w, http.StatusConflict, err.Error(), nil)
		case contracts.ErrInvalidCredentials:
			writeError(w, http.StatusUnauthorized, err.Error(), nil)
		default:
			writeError(w, http.StatusInternalServerError, "internal server error", nil)
		}
	}
}
