package seeder

import (
	"github.com/arman300s/uni-portal/internal/models"
	"gorm.io/gorm"
)

func SeedRoles(db *gorm.DB) {
	roles := []models.Role{
		{Name: "admin"},
		{Name: "teacher"},
		{Name: "student"},
	}

	for _, role := range roles {
		db.FirstOrCreate(&role, models.Role{Name: role.Name})
	}
}
