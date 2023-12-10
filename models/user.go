package models

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Username  string `gorm:"uniqueIndex;not null;size:50;" validate:"required,min=3,max=50" json:"username"`
	Email     string `gorm:"uniqueIndex;not null;size:255;" validate:"required,email" json:"email"`
	Password  string `gorm:"not null;" validate:"required,min=6,max=50" json:"password"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	IsAdmin   bool   `gorm:"default:false" json:"is_admin"`
}
