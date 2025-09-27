package models

import "time"

type User struct {
	ID        uint      `gorm:"primary_key"`
	Name      string    `gorm:"size:100;not null"`
	Email     string    `gorm:"UniqueIndex;not null"`
	Password  string    `gorm:"size:255;not null"`
	CreatedAt time.Time `gorm:"DEFAULT:CURRENT_TIMESTAMP"`
	UpdatedAt time.Time `gorm:"DEFAULT:CURRENT_TIMESTAMP"`
}
