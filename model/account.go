package model

import (
	"gorm.io/gorm"
)

type Account struct {
	gorm.Model
	StakeKey        string            `gorm:"not null" json:"stakeKey"`
	Title           string            `json:"title"`
	Users           []*User           `gorm:"many2many:user_accounts;"`
	StakeKeyHistory []*AccountHistory `gorm:"foreignKey:AccountID;references:ID"`
	Transactions    []*Transaction    `gorm:"foreignKey:AccountID;references:ID"`
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

// get user accounts details
func (A *Account) GetUserAccountsDetails(db *gorm.DB, userID uint) ([]Account, error) {
	var accounts []Account
	result := db.Debug().Where("user_id = ?", userID).Preload("UserAccounts").Find(&accounts)
	if result.Error != nil {
		return accounts, result.Error
	}
	return accounts, nil
}

