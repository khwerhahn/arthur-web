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
		go func() {
			for {
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

				// send to channel
				streamChan <- returnString
				time.Sleep(time.Second * 10)
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
