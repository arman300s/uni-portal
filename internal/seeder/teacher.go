package seeder

import (
	"log"

	"github.com/arman300s/uni-portal/internal/models"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func SeedTeachers(database *gorm.DB) {
	var teacherRole models.Role
	if err := database.First(&teacherRole, "name = ?", "teacher").Error; err != nil {
		log.Println("⚠️ Teacher role not found, skipping teacher seeding")
		return
	}

	teachers := []models.User{
		{Name: "Mr Avinash", Email: "avinash@uni.kz"},
		{Name: "Murat Abdilda", Email: "abdilda@uni.kz"},
		{Name: "Torekeldi Niyazbek", Email: "torekeldi@uni.kz"},
	}

	for _, t := range teachers {
		var existing models.User
		if err := database.First(&existing, "email = ?", t.Email).Error; err == nil {
			continue
		}

		password, _ := bcrypt.GenerateFromPassword([]byte("123456"), bcrypt.DefaultCost)
		t.Password = string(password)
		t.RoleID = &teacherRole.ID

		if err := database.Create(&t).Error; err != nil {
			log.Printf("❌ Failed to seed teacher %s: %v\n", t.Email, err)
		}
	}
}
