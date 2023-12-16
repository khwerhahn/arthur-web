package model

import (
	"gorm.io/gorm"
)

type StakeKeyHistory struct {
	gorm.Model
	AccountID        uint    `gorm:"not null" json:"accountId"`
	EpochID          uint    `gorm:"not null" json:"epochId"`
	EpochIdAvailable uint    `gorm:"not null" json:"epochIdAvailable"`
	Amount           float64 `gorm:"not null" json:"amount"`
	ControlledAmount float64 `gorm:"not null;default:0" json:"controlledAmount"`
	Pool             string  `gorm:"not null" json:"pool"`
	Type             string  `gorm:"not null" json:"type"`
}

func (StakeKeyHistory) TableName() string {
	return "stakekeyhistory"
}

// func is for stake key sync to just get last epoch the stake key was synced at
func (StakeKeyHistory) GetCheckStakeKeyHistory(db *gorm.DB, accountID uint) ([]StakeKeyHistory, error) {
	// raw query
	var response []StakeKeyHistory
	err := db.Raw("SELECT s.id,s.account_id,e.id AS epoch_id,e.title,s.epoch_id_available,s.amount,s.pool,s.type FROM stakekeyhistory s,epochs e WHERE s.epoch_id=e.ID AND account_id= ? ORDER BY e.title DESC LIMIT 1", accountID).Scan(&response).Error
	if err != nil {
		return response, err
	}
	return response, nil
}
