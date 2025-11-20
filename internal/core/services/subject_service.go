package services

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"gorm.io/gorm"

	"github.com/arman300s/uni-portal/internal/core/contracts"
	"github.com/arman300s/uni-portal/internal/core/repositories"
	"github.com/arman300s/uni-portal/internal/models"
)

// SubjectService manages subject-related logic.
type SubjectService struct {
	subjects repositories.SubjectRepository
	users    repositories.UserRepository
}

func NewSubjectService(subjects repositories.SubjectRepository, users repositories.UserRepository) *SubjectService {
	return &SubjectService{subjects: subjects, users: users}
}

func (s *SubjectService) CreateSubject(ctx context.Context, input contracts.SubjectInput) (*models.Subject, error) {
	if errs := validateSubjectInput(input); len(errs) > 0 {
		return nil, errs
	}

	subject := &models.Subject{
		Name:        strings.TrimSpace(input.Name),
		Description: strings.TrimSpace(input.Description),
	}

	if len(input.TeacherIDs) > 0 {
		teachers, err := s.fetchTeacherUsers(ctx, input.TeacherIDs)
		if err != nil {
			return nil, err
		}
		subject.Teachers = teachers
	}

	if err := s.subjects.Create(ctx, subject); err != nil {
		return nil, err
	}

	return subject, nil
}

func (s *SubjectService) ListSubjects(ctx context.Context) ([]models.Subject, error) {
	return s.subjects.List(ctx)
}

func (s *SubjectService) GetSubject(ctx context.Context, id uint) (*models.Subject, error) {
	subject, err := s.subjects.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, contracts.ErrSubjectNotFound
		}
		return nil, err
	}
	return subject, nil
}

func (s *SubjectService) UpdateSubject(ctx context.Context, id uint, input contracts.SubjectInput) error {
	subject, err := s.subjects.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return contracts.ErrSubjectNotFound
		}
		return err
	}

	if trimmed := strings.TrimSpace(input.Name); trimmed != "" {
		subject.Name = trimmed
	}
	if desc := strings.TrimSpace(input.Description); desc != "" {
		subject.Description = desc
	}

	if len(input.TeacherIDs) > 0 {
		teachers, err := s.fetchTeacherUsers(ctx, input.TeacherIDs)
		if err != nil {
			return err
		}
		if err := s.subjects.ReplaceTeachers(ctx, subject, teachers); err != nil {
			return err
		}
		subject.Teachers = teachers
	}

	return s.subjects.Save(ctx, subject)
}

func (s *SubjectService) DeleteSubject(ctx context.Context, id uint) error {
	if err := s.subjects.Delete(ctx, id); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return contracts.ErrSubjectNotFound
		}
		return err
	}
	return nil
}

func (s *SubjectService) ListSubjectsForTeacher(ctx context.Context, teacherID uint) ([]models.Subject, error) {
	return s.subjects.ListByTeacherID(ctx, teacherID)
}

func (s *SubjectService) fetchTeacherUsers(ctx context.Context, ids []uint) ([]models.User, error) {
	users, err := s.users.FindByIDs(ctx, ids)
	if err != nil {
		return nil, err
	}

	if len(users) != len(ids) {
		return nil, contracts.ValidationErrors{contracts.ValidationError{
			Field:   "teacher_ids",
			Message: "one or more teachers not found",
		}}
	}

	idLookup := make(map[uint]struct{}, len(ids))
	for _, id := range ids {
		idLookup[id] = struct{}{}
	}

	for _, u := range users {
		if _, ok := idLookup[u.ID]; !ok {
			return nil, contracts.ValidationErrors{contracts.ValidationError{
				Field:   "teacher_ids",
				Message: fmt.Sprintf("invalid teacher id %d", u.ID),
			}}
		}
		if u.Role == nil || u.Role.Name != "teacher" {
			return nil, contracts.ValidationErrors{contracts.ValidationError{
				Field:   "teacher_ids",
				Message: fmt.Sprintf("user %d is not a teacher", u.ID),
			}}
		}
	}

	return users, nil
}
