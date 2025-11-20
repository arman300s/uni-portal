package controllers

import (
	"net/http"

	"github.com/arman300s/uni-portal/internal/core/contracts"
	"github.com/arman300s/uni-portal/internal/core/services"
	"github.com/arman300s/uni-portal/pkg/middleware"
)

// TeacherController exposes teacher-specific endpoints.
type TeacherController struct {
	service *services.SubjectService
}

func NewTeacherController(service *services.SubjectService) *TeacherController {
	return &TeacherController{service: service}
}

// ListMySubjects godoc
// @Summary List teacher subjects
// @Tags teacher-subjects
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {array} contracts.SubjectDTO
// @Failure 401 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /teacher/subjects [get]
func (c *TeacherController) ListMySubjects(w http.ResponseWriter, r *http.Request) {
	teacherID, ok := middleware.UserIDFromContext(r.Context())
	if !ok {
		writeError(w, http.StatusUnauthorized, "unauthorized", nil)
		return
	}

	subjects, err := c.service.ListSubjectsForTeacher(r.Context(), teacherID)
	if err != nil {
		handleSubjectError(w, err)
		return
	}

	resp := make([]contracts.SubjectDTO, 0, len(subjects))
	for _, s := range subjects {
		resp = append(resp, contracts.SubjectDTO{
			ID:          s.ID,
			Name:        s.Name,
			Description: s.Description,
			Teachers:    extractTeacherNames(s.Teachers),
		})
	}

	writeJSON(w, http.StatusOK, resp)
}
