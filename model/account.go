package model

import (
	"gorm.io/gorm"
)

type Account struct {
	gorm.Model
	StakeKey        string             `gorm:"not null" json:"stakeKey"`
	Title           string             `json:"title"`
	Users           []*User            `gorm:"many2many:users_accounts;"`
	StakeKeyHistory []*StakeKeyHistory `gorm:"foreignKey:AccountID;references:ID"`
}

func (A *Account) GetAccountByStakeKey(db *gorm.DB, stakeKey string) error {
	result := db.Where("stake_key = ?", stakeKey).First(&A)
	if result.Error != nil {
		return result.Error
	}
	return nil
}
