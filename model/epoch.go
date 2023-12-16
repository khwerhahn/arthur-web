package model

import (
	"time"

	"gorm.io/gorm"
)

type Epoch struct {
	gorm.Model
	Title      uint       `gorm:"not null;type:int4" json:"title"`
	EpochStart *time.Time `gorm:"not null" json:"epochStart"`
	EpochEnd   *time.Time `gorm:"not null" json:"epochEnd"`
}

func (Epoch) TableName() string {
	return "epochs"
}

func (E *Epoch) GetCurrentEpoch(db *gorm.DB) error {
	// get current epoch where epoch start is smaller or equal to now and epoch end is greater or equal to now
	result := db.Where("epoch_start <= ? AND epoch_end >= ?", time.Now(), time.Now()).First(&E)
	if result.Error != nil {
		return result.Error
	}
	return nil
}
