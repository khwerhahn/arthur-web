package model

import "gorm.io/gorm"

type UserAccounts struct {
	gorm.Model
	UserID    uint   `gorm:"not null;index:idx_user_account,unique" json:"userId"`
	AccountID uint   `gorm:"not null;index:idx_user_account,unique" json:"accountId"`
	Title     string `json:"title"`
}

func (u *UserAccounts) SaveUserAccount(db *gorm.DB) (uint, error) {
	result := db.Create(&u)
	if result.Error != nil {
		return 0, result.Error
	}
	return u.ID, nil
}

func (u *UserAccounts) GetUserAccounts(db *gorm.DB, userID uint) ([]UserAccounts, error) {
	var userAccounts []UserAccounts
	result := db.Where("user_id = ?", userID).Find(&userAccounts)
	if result.Error != nil {
		return userAccounts, result.Error
	}
	return userAccounts, nil
}
