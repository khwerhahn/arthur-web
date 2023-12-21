package model

import (
	"gorm.io/gorm"
)

type Account struct {
	gorm.Model
	StakeKey        string             `gorm:"not null;unique" json:"stakeKey"`
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

// save
func (A *Account) SaveAccount(db *gorm.DB) (*Account, error) {
	result := db.Create(&A)
	if result.Error != nil {
		return A, result.Error
	}
	return A, nil
}

// get account by id
func (A *Account) GetAccountByID(db *gorm.DB, id uint) (Account, error) {
	result := db.Where("id = ?", id).Find(&A)
	if result.Error != nil {
		return *A, result.Error
	}
	return *A, nil
}

