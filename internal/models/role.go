package models

import "time"

type Role struct {
	ID        uint   `gorm:"primary_key"`
	Name      string `gorm:"size:50;unique;not null"`
	CreatedAt time.Time
	UpdatedAt time.Time
}
