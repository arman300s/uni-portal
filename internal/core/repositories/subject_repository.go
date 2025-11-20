package repositories

import (
	"context"

	"github.com/arman300s/uni-portal/internal/models"
	"gorm.io/gorm"
)

// SubjectRepository exposes persistence operations for subjects.
type SubjectRepository interface {
	Create(ctx context.Context, subject *models.Subject) error
	List(ctx context.Context) ([]models.Subject, error)
	FindByID(ctx context.Context, id uint) (*models.Subject, error)
	Save(ctx context.Context, subject *models.Subject) error
	Delete(ctx context.Context, id uint) error
	ReplaceTeachers(ctx context.Context, subject *models.Subject, teachers []models.User) error
	ListByTeacherID(ctx context.Context, teacherID uint) ([]models.Subject, error)
}

type subjectRepository struct {
	db *gorm.DB
}

func NewSubjectRepository(db *gorm.DB) SubjectRepository {
	return &subjectRepository{db: db}
}

func (r *subjectRepository) Create(ctx context.Context, subject *models.Subject) error {
	return r.db.WithContext(ctx).Create(subject).Error
}

func (r *subjectRepository) List(ctx context.Context) ([]models.Subject, error) {
	var subjects []models.Subject
	if err := r.db.WithContext(ctx).Preload("Teachers").Find(&subjects).Error; err != nil {
		return nil, err
	}
	return subjects, nil
}

func (r *subjectRepository) FindByID(ctx context.Context, id uint) (*models.Subject, error) {
	var subject models.Subject
	if err := r.db.WithContext(ctx).Preload("Teachers").First(&subject, id).Error; err != nil {
		return nil, err
	}
	return &subject, nil
}

func (r *subjectRepository) Save(ctx context.Context, subject *models.Subject) error {
	return r.db.WithContext(ctx).Save(subject).Error
}

func (r *subjectRepository) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&models.Subject{}, id).Error
}

func (r *subjectRepository) ReplaceTeachers(ctx context.Context, subject *models.Subject, teachers []models.User) error {
	return r.db.WithContext(ctx).Model(subject).Association("Teachers").Replace(teachers)
}

func (r *subjectRepository) ListByTeacherID(ctx context.Context, teacherID uint) ([]models.Subject, error) {
	var subjects []models.Subject
	if err := r.db.WithContext(ctx).
		Preload("Teachers").
		Joins("JOIN subject_teachers st ON st.subject_id = subjects.id").
		Where("st.user_id = ?", teacherID).
		Find(&subjects).Error; err != nil {
		return nil, err
	}
	return subjects, nil
}
