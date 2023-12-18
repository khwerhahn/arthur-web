package model

import (
	"time"

	"gorm.io/gorm"
)

type MarketData struct {
	gorm.Model
	Name              string     `gorm:"not null" json:"name"`
	Symbol            string     `gorm:"not null;index:idx_comp,unique" json:"symbol"`
	DataProvider      string     `gorm:"not null;default:cg" json:"dataProvider"`
	QuoteDenomination string     `gorm:"not null;index:idx_comp,unique" json:"quoteDenomination"`
	Open              float64    `gorm:"not null" json:"open"`
	Close             float64    `gorm:"not null" json:"close"`
	High              float64    `gorm:"not null" json:"high"`
	Low               float64    `gorm:"not null" json:"low"`
	Volume            float64    `json:"volume"`
	RangeFrom         *time.Time `gorm:"not null;index:idx_comp,unique" json:"rangeFrom"`
	RangeTo           *time.Time `gorm:"not null;index:idx_comp,unique" json:"rangeTo"`
	TimestampRemote   int64      `json:"timestampRemote"`
	RemoteRetry       bool       `gorm:"not null;default:true" json:"remoteRetry"`
	UpdateCount       int        `gorm:"not null;default:0" json:"updateCount"`
}

func (MarketData) TableName() string {
	return "marketdata"
}

// get marketdata between two timestamps
func (M *MarketData) GetMarketDataBetween(db *gorm.DB, symbol string, quoteDenomination string, rangeFrom *time.Time, rangeTo *time.Time) ([]MarketData, error) {
	var marketData []MarketData
	result := db.Where("symbol = ? AND quote_denomination = ? AND range_from >= ? AND range_to <= ?", symbol, quoteDenomination, rangeFrom, rangeTo).Find(&marketData)
	if result.Error != nil {
		return []MarketData{}, result.Error
	}
	return marketData, nil
}

// check if marketdata exists
func (M *MarketData) Exists(db *gorm.DB) bool {
	var marketData MarketData
	result := db.Where("timestamp_remote = ? AND quote_denomination = ? AND symbol = ?", M.TimestampRemote, M.QuoteDenomination, M.Symbol).First(&marketData)
	if result.Error != nil {
		return false
	}
	return true
}

func (M *MarketData) Save(db *gorm.DB) error {
	result := db.Save(&M)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

// hook for create to increment updateCount
func (M *MarketData) BeforeCreate(tx *gorm.DB) (err error) {
	M.UpdateCount++
	return
}

// hook for update to increment UpdateCount
func (M *MarketData) BeforeUpdate(tx *gorm.DB) (err error) {
	M.UpdateCount++
	return
}

// get last availeble marketdata
func (M *MarketData) GetLastAvailableMarketData(db *gorm.DB, symbol string, quoteDenomination string) (MarketData, error) {
	var marketData MarketData
	result := db.Where("symbol = ? AND quote_denomination = ? ", symbol, quoteDenomination).Order("range_from desc").First(&marketData)
	if result.Error != nil {
		return MarketData{}, result.Error
	}
	return marketData, nil
}
