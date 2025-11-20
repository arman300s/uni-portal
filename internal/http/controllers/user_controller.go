package controllers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"

	"github.com/arman300s/uni-portal/internal/core/contracts"
	"github.com/arman300s/uni-portal/internal/core/services"
	"github.com/arman300s/uni-portal/pkg/middleware"
)

// UserController manages user/admin endpoints.
type UserController struct {
	service *services.UserService
}

func NewUserController(service *services.UserService) *UserController {
	return &UserController{service: service}
}

// Me godoc
// @Summary Get current user
// @Description Get current authenticated user information
// @Tags auth
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {object} contracts.UserDTO
// @Failure 401 {object} ErrorResponse
// @Router /me [get]
func (c *UserController) Me(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.UserIDFromContext(r.Context())
	if !ok {
		writeError(w, http.StatusUnauthorized, "unauthorized", nil)
		return
	}

	user, err := c.service.GetCurrentUser(r.Context(), userID)
	if err != nil {
		handleUserError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, user)
}

// ListUsers godoc
// @Summary List users
// @Tags admin-users
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {array} contracts.UserDTO
// @Failure 500 {object} ErrorResponse
// @Router /admin/users [get]
func (c *UserController) ListUsers(w http.ResponseWriter, r *http.Request) {
	users, err := c.service.ListUsers(r.Context())
	if err != nil {
		handleUserError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, users)
}

// GetUser godoc
// @Summary Get user
// @Tags admin-users
// @Produce json
// @Security ApiKeyAuth
// @Param id path int true "User ID"
// @Success 200 {object} contracts.UserDTO
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Router /admin/users/{id} [get]
func (c *UserController) GetUser(w http.ResponseWriter, r *http.Request) {
	id, err := parseIDParam(r)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error(), nil)
		return
	}

	user, err := c.service.GetUser(r.Context(), id)
	if err != nil {
		handleUserError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, user)
}

// CreateUser godoc
// @Summary Create user
// @Tags admin-users
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param user body contracts.CreateUserInput true "User payload"
// @Success 201 {object} map[string]interface{}
// @Failure 400 {object} ErrorResponse
// @Failure 409 {object} ErrorResponse
// @Router /admin/users/create [post]
func (c *UserController) CreateUser(w http.ResponseWriter, r *http.Request) {
	var input contracts.CreateUserInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body", nil)
		return
	}

	user, err := c.service.CreateUser(r.Context(), input)
	if err != nil {
		handleUserError(w, err)
		return
	}

	writeJSON(w, http.StatusCreated, map[string]interface{}{
		"message": "user created successfully",
		"user":    user,
	})
}

// UpdateUser godoc
// @Summary Update user role
// @Tags admin-users
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param id path int true "User ID"
// @Param user body contracts.UpdateUserInput true "Update payload"
// @Success 200 {object} map[string]string
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Router /admin/users/{id} [put]
func (c *UserController) UpdateUser(w http.ResponseWriter, r *http.Request) {
	id, err := parseIDParam(r)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error(), nil)
		return
	}

	var input contracts.UpdateUserInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body", nil)
		return
	}

	if err := c.service.UpdateUserRole(r.Context(), id, input); err != nil {
		handleUserError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{"message": "user updated successfully"})
}

// DeleteUser godoc
// @Summary Delete user
// @Tags admin-users
// @Produce json
// @Security ApiKeyAuth
// @Param id path int true "User ID"
// @Success 200 {object} map[string]string
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Router /admin/users/{id} [delete]
func (c *UserController) DeleteUser(w http.ResponseWriter, r *http.Request) {
	id, err := parseIDParam(r)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error(), nil)
		return
	}

	if err := c.service.DeleteUser(r.Context(), id); err != nil {
		handleUserError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{"message": "user deleted successfully"})
}

func parseIDParam(r *http.Request) (uint, error) {
	idStr := mux.Vars(r)["id"]
	id64, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		return 0, err
	}
	return uint(id64), nil
}

func handleUserError(w http.ResponseWriter, err error) {
	switch e := err.(type) {
	case contracts.ValidationErrors:
		writeError(w, http.StatusBadRequest, "validation failed", e)
		return
	}

	switch err {
	case contracts.ErrUserNotFound:
		writeError(w, http.StatusNotFound, err.Error(), nil)
	case contracts.ErrRoleNotFound:
		writeError(w, http.StatusBadRequest, err.Error(), nil)
	case contracts.ErrEmailInUse:
		writeError(w, http.StatusConflict, err.Error(), nil)
	default:
		writeError(w, http.StatusInternalServerError, "internal server error", nil)
	}
}
