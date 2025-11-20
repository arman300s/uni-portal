package controllers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"

	"github.com/arman300s/uni-portal/internal/core/contracts"
	"github.com/arman300s/uni-portal/internal/core/services"
	"github.com/arman300s/uni-portal/internal/models"
)

// AdminSubjectController manages subject CRUD for admins.
type AdminSubjectController struct {
	service *services.SubjectService
}

func NewAdminSubjectController(service *services.SubjectService) *AdminSubjectController {
	return &AdminSubjectController{service: service}
}

// CreateSubject godoc
// @Summary Create subject
// @Tags admin-subjects
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param subject body contracts.SubjectInput true "Subject payload"
// @Success 201 {object} map[string]interface{}
// @Failure 400 {object} ErrorResponse
// @Router /admin/subjects [post]
func (c *AdminSubjectController) CreateSubject(w http.ResponseWriter, r *http.Request) {
	var input contracts.SubjectInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body", nil)
		return
	}

	subject, err := c.service.CreateSubject(r.Context(), input)
	if err != nil {
		handleSubjectError(w, err)
		return
	}

	writeJSON(w, http.StatusCreated, map[string]interface{}{
		"message": "subject created successfully",
		"subject": subjectResponse(subject),
	})
}

// ListSubjects godoc
// @Summary List subjects
// @Tags admin-subjects
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {array} map[string]interface{}
// @Failure 500 {object} ErrorResponse
// @Router /admin/subjects [get]
func (c *AdminSubjectController) ListSubjects(w http.ResponseWriter, r *http.Request) {
	subjects, err := c.service.ListSubjects(r.Context())
	if err != nil {
		handleSubjectError(w, err)
		return
	}

	payload := make([]map[string]interface{}, 0, len(subjects))
	for i := range subjects {
		payload = append(payload, subjectResponse(&subjects[i]))
	}

	writeJSON(w, http.StatusOK, payload)
}

// GetSubject godoc
// @Summary Get subject
// @Tags admin-subjects
// @Produce json
// @Security ApiKeyAuth
// @Param id path int true "Subject ID"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Router /admin/subjects/{id} [get]
func (c *AdminSubjectController) GetSubject(w http.ResponseWriter, r *http.Request) {
	id, err := parseSubjectID(r)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error(), nil)
		return
	}

	subject, err := c.service.GetSubject(r.Context(), id)
	if err != nil {
		handleSubjectError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, subjectResponse(subject))
}

// UpdateSubject godoc
// @Summary Update subject
// @Tags admin-subjects
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param id path int true "Subject ID"
// @Param subject body contracts.SubjectInput true "Subject payload"
// @Success 200 {object} map[string]string
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Router /admin/subjects/{id} [put]
func (c *AdminSubjectController) UpdateSubject(w http.ResponseWriter, r *http.Request) {
	id, err := parseSubjectID(r)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error(), nil)
		return
	}

	var input contracts.SubjectInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body", nil)
		return
	}

	if err := c.service.UpdateSubject(r.Context(), id, input); err != nil {
		handleSubjectError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{"message": "subject updated successfully"})
}

// DeleteSubject godoc
// @Summary Delete subject
// @Tags admin-subjects
// @Produce json
// @Security ApiKeyAuth
// @Param id path int true "Subject ID"
// @Success 200 {object} map[string]string
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Router /admin/subjects/{id} [delete]
func (c *AdminSubjectController) DeleteSubject(w http.ResponseWriter, r *http.Request) {
	id, err := parseSubjectID(r)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error(), nil)
		return
	}

	if err := c.service.DeleteSubject(r.Context(), id); err != nil {
		handleSubjectError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{"message": "subject deleted successfully"})
}

func parseSubjectID(r *http.Request) (uint, error) {
	idStr := mux.Vars(r)["id"]
	val, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		return 0, err
	}
	return uint(val), nil
}

func subjectResponse(subject *models.Subject) map[string]interface{} {
	teachers := make([]map[string]interface{}, 0, len(subject.Teachers))
	for _, t := range subject.Teachers {
		teachers = append(teachers, map[string]interface{}{
			"id":    t.ID,
			"name":  t.Name,
			"email": t.Email,
		})
	}
	return map[string]interface{}{
		"id":          subject.ID,
		"name":        subject.Name,
		"description": subject.Description,
		"teachers":    teachers,
	}
}

func handleSubjectError(w http.ResponseWriter, err error) {
	switch e := err.(type) {
	case contracts.ValidationErrors:
		writeError(w, http.StatusBadRequest, "validation failed", e)
		return
	}

	switch err {
	case contracts.ErrSubjectNotFound:
		writeError(w, http.StatusNotFound, err.Error(), nil)
	default:
		writeError(w, http.StatusInternalServerError, "internal server error", nil)
	}
}
