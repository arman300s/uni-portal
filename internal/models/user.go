package models

import "time"

type User struct {
	ID        uint   `gorm:"primary_key"`
	Name      string `gorm:"size:100;not null"`
	Email     string `gorm:"UniqueIndex;not null"`
	Password  string `gorm:"size:255;not null"`
	RoleID    *uint  `gorm:"default:null"`
	Role      *Role
	CreatedAt time.Time `gorm:"DEFAULT:CURRENT_TIMESTAMP"`
	UpdatedAt time.Time `gorm:"DEFAULT:CURRENT_TIMESTAMP"`
}
