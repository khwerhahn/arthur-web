package model

import "gorm.io/gorm"

type UsersAccounts struct {
	gorm.Model
	UserID    uint   `gorm:"not null;foreignKey:UserID;index:idx_user_account,unique" json:"userId"`
	AccountID uint   `gorm:"not null;foreignKey:AccountID;index:idx_user_account,unique" json:"accountId"`
	Title     string `json:"title"`
}

func (u *UsersAccounts) SaveUserAccount(db *gorm.DB) (uint, error) {
	result := db.Create(&u)
	if result.Error != nil {
		return 0, result.Error
	}
	return u.ID, nil
}
