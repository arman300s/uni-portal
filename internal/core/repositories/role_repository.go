package repositories

import (
	"context"

	"github.com/arman300s/uni-portal/internal/models"
	"gorm.io/gorm"
)

type RoleRepository interface {
	FindByName(ctx context.Context, name string) (*models.Role, error)
}

type roleRepository struct {
	db *gorm.DB
}

func NewRoleRepository(db *gorm.DB) RoleRepository {
	return &roleRepository{db: db}
}

func (r *roleRepository) FindByName(ctx context.Context, name string) (*models.Role, error) {
	var role models.Role
	if err := r.db.WithContext(ctx).Where("name = ?", name).First(&role).Error; err != nil {
		return nil, err
	}
	return &role, nil
}
