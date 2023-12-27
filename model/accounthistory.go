package model

import (
	"gorm.io/gorm"
)

type AccountHistory struct {
	gorm.Model
	AccountID        uint   `gorm:"not null;index:idx_unhis,unique" json:"accountId"`
	EpochID          uint   `gorm:"not null;index:idx_unhis,unique" json:"epochId"`
	EpochIdAvailable uint   `gorm:"not null" json:"epochIdAvailable"`
	RewardsAmount    int64  `gorm:"not null;default:0" json:"rewardsAmount"`
	Amount           int64  `gorm:"not null;default:0" json:"amount"`
	Pool             string `gorm:"not null" json:"pool"`
	Type             string `gorm:"not null" json:"type"`
	Epoch            Epoch  `gorm:"foreignKey:EpochID;references:ID"`
}

func (AccountHistory) TableName() string {
	return "accounthistory"
}

// func is for stake key sync to just get last epoch the stake key was synced at
func (AccountHistory) GetCheckStakeKeyHistory(db *gorm.DB, accountID uint) ([]AccountHistory, error) {
	// raw query
	var response []AccountHistory
	err := db.Raw("SELECT s.id,s.account_id,e.id AS epoch_id,e.title,s.epoch_id_available,s.amount,s.pool,s.type FROM stakekeyhistory s,epochs e WHERE s.epoch_id=e.ID AND account_id= ? ORDER BY e.title DESC LIMIT 1", accountID).Scan(&response).Error
	if err != nil {
		return response, err
	}
	return response, nil
}
