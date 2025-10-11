package seeder

import (
	"log"

	"github.com/arman300s/uni-portal/internal/models"
	"gorm.io/gorm"
)

func SeedSubjects(database *gorm.DB) {
	subjects := []models.Subject{
		{Name: "Mathematics"},
		{Name: "Computer Science"},
		{Name: "Physics"},
		{Name: "English Language"},
		{Name: "History"},
	}

	for _, s := range subjects {
		var existing models.Subject
		if err := database.First(&existing, "name = ?", s.Name).Error; err == nil {
			continue
		}

		if err := database.Create(&s).Error; err != nil {
			log.Printf("❌ Failed to seed subject %s: %v\n", s.Name, err)
		}
	}
}
