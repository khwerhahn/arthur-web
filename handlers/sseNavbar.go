package handlers

import (
	"arthur-web/globals"
	"arthur-web/model"
	"io"
	"strconv"
	"time"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func SseNavbarContents(DB *gorm.DB, currencyToUse string) string {
	returnString := ""
	// get last available epoch
	var marketData model.MarketData
	lastMarketData, err := marketData.GetLastAvailableMarketData(DB, "ada", currencyToUse)
	if err != nil {
		returnString = "error"
	} else {
		// return
		// Close string + currneycurrencyToUse string
		currSymbol := "$"
		if currencyToUse == "eur" {
			currSymbol = "€"
		}
		returnString = strconv.FormatFloat(lastMarketData.Close, 'f', 2, 64) + " " + currSymbol
	}
	return returnString
}

func SseNavbar(DB *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		currencyToUse := "usd"
		// get session
		session := sessions.Default(c)
		// get user user setting currency
		userCurrency := session.Get(globals.UserSettingCurrency)
		if userCurrency != nil {
			currencyToUse = userCurrency.(string)
		}

		streamChan := make(chan string)

		// ticker
		quit := make(chan struct{})
		ticker := time.NewTicker(time.Second * 10)
		// just the inital chan 1 microsecond
		tickerOnceOnly := time.NewTicker(time.Microsecond * 1)

		go func() {
			for {
				select {
				case <-ticker.C:
					returnString := SseNavbarContents(DB, currencyToUse)
					streamChan <- returnString
				case <-tickerOnceOnly.C:
					returnString := SseNavbarContents(DB, currencyToUse)
					streamChan <- returnString
					tickerOnceOnly.Stop()
				case <-quit:
					ticker.Stop()
					return
				}
			}
		}()

		c.Stream(func(w io.Writer) bool {
			// Stream message to client from message channel
			if msg, ok := <-streamChan; ok {
				c.SSEvent("message", msg)
				return true
			}
			return false
		})
	}
}
