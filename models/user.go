package models

import (
	"html"
	"strings"
	"time"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// the model Session is fed by this struct

type User struct {
	gorm.Model
	Username  string     `gorm:"unique;not null;size:50;" validate:"required,min=3,max=50" json:"username"`
	Email     string     `gorm:"unique;not null;size:255;" validate:"required,email" json:"email"`
	Password  string     `gorm:"not null;" validate:"required,min=6,max=50" json:"password"`
	FirstName string     `gorm:"not null;size:50;" validate:"required,min=3,max=50" json:"first_name"`
	LastName  string     `gorm:"not null;size:50;" validate:"required,min=3,max=200" json:"last_name"`
	IsAdmin   bool       `gorm:"default:false" json:"is_admin"`
	CreatedAt *time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt *time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
}

func (u *User) SaveUser(DB *gorm.DB) (*User, error) {
	var err error
	err = DB.Create(&u).Error
	if err != nil {
		return &User{}, err
	}
	return u, nil
}

func (u *User) BeforeCreate(tx *gorm.DB) error {

	//turn password into hash
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.Password = string(hashedPassword)

	//remove spaces in username
	u.Username = html.EscapeString(strings.TrimSpace(u.Username))

	return nil

}
