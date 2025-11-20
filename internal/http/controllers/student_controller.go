package controllers

import (
    "net/http"

    "github.com/arman300s/uni-portal/internal/core/contracts"
    "github.com/arman300s/uni-portal/internal/core/services"
    "github.com/arman300s/uni-portal/internal/models"
)

// StudentController exposes student endpoints.
type StudentController struct {
    service *services.SubjectService
}

func NewStudentController(service *services.SubjectService) *StudentController {
    return &StudentController{service: service}
}

// ListSubjects godoc
// @Summary List subjects for students
// @Tags student-subjects
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {array} contracts.SubjectDTO
// @Failure 500 {object} ErrorResponse
// @Router /student/subjects [get]
func (c *StudentController) ListSubjects(w http.ResponseWriter, r *http.Request) {
    subjects, err := c.service.ListSubjects(r.Context())
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

func extractTeacherNames(teachers []models.User) []string {
    names := make([]string, 0, len(teachers))
    for _, t := range teachers {
        names = append(names, t.Name)
    }
    return names
}
