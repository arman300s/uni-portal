package models

import "gorm.io/gorm"

type Subject struct {
	gorm.Model
	Name        string `json:"name" gorm:"unique;not null"`
	Description string `json:"description"`

	Teachers []User `json:"teachers" gorm:"many2many:subject_teachers;constraint:OnDelete:CASCADE;"`
}
