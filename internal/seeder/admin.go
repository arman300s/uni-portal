package seeder

import (
	"log"

	"github.com/arman300s/uni-portal/internal/models"
	"github.com/arman300s/uni-portal/pkg/auth" // for password hashing
	"gorm.io/gorm"
)

func SeedAdmin(db *gorm.DB) {
	SeedRoles(db)

	var adminRole models.Role
	if err := db.Where("name = ?", "admin").First(&adminRole).Error; err != nil {
		log.Printf("❌ Failed to find admin role: %v\n", err)
		return
	}

	var adminUser models.User
	if err := db.Where("email = ?", "admin@uni-portal.com").First(&adminUser).Error; err == nil {
		log.Println("ℹ️ Admin user already exists")
		return
	}

	hash, err := auth.HashPassword("admin123")
	if err != nil {
		log.Printf("❌ Failed to hash admin password: %v\n", err)
		return
	}

	adminUser = models.User{
		Name:     "Super Admin",
		Email:    "admin@uni-portal.com",
		Password: hash,
		RoleID:   &adminRole.ID,
	}

	if err := db.Create(&adminUser).Error; err != nil {
		log.Printf("❌ Failed to create admin user: %v\n", err)
		return
	}

	log.Println("✅ Default admin user created successfully (admin@uni-portal.com / admin123)")
}
