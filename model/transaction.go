package model

import "gorm.io/gorm"

type Transaction struct {
	gorm.Model
	AccountID uint   `gorm:"not null" json:"accountID"`
	Hash      string `gorm:"not null" json:"hash"`
}
