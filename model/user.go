package model

import (
	"encoding/json"
	"html"
	"strings"
	"time"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// the model Session is fed by this struct
type UserSettings struct {
	Currency string `json:"currency"`
	Language string `json:"language"`
}

type User struct {
	gorm.Model
	Username        string     `gorm:"unique;not null;size:50;" validate:"required,min=3,max=50" json:"username"`
	Email           string     `gorm:"unique;not null;size:255;" validate:"required,email" json:"email"`
	Password        string     `gorm:"not null;" validate:"required,min=6,max=50" json:"password"`
	FirstName       string     `gorm:"not null;size:50;" validate:"required,min=3,max=50" json:"first_name"`
	LastName        string     `gorm:"not null;size:50;" validate:"required,min=3,max=200" json:"last_name"`
	IsAdmin         bool       `gorm:"default:false" json:"is_admin"`
	ProfileImageUrl string     `gorm:"default:'https://i.pravatar.cc/100?img=15'" json:"image_url"`
	UserSettings string     `gorm:"not null;default: '{\"currency\": \"usd\", \"language\": \"en\"}'" json:"userSettings"`
	Accounts     []*Account `gorm:"many2many:users_accounts;"`
	CreatedAt       *time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt       *time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
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

func (u *User) ComparePassword(password string) error {
	return bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
}

// Get User Seetings
func (u *User) GetUserSettings() (UserSettings, error) {
	// decode json from database to struct
	var userSettings UserSettings
	err := json.Unmarshal([]byte(u.UserSettings), &userSettings)
	if err != nil {
		return UserSettings{}, err
	}
	return userSettings, nil
}

// Set User UserSettings
func (u *User) SetUserSettings(userSettings UserSettings) error {
	// encode struct to json for database
	userSettingsJson, err := json.Marshal(userSettings)
	if err != nil {
		return err
	}
	u.UserSettings = string(userSettingsJson)
	return nil
}

// get user by id and return user
func (u *User) GetUserByID(db *gorm.DB, id uint) (*User, error) {
	result := db.Where("id = ?", id).First(&u)
	if result.Error != nil {
		return u, result.Error
	}
	return u, nil
}
