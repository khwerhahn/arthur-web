package handlers

import (
	"arthur-web/globals"
	"arthur-web/views"
	"bytes"
	"context"
	"fmt"
	"io"
	"time"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func SseWalletContents(DB *gorm.DB, currencyToUse string, walletID string, activeUrl string, action string) string {

	// use templ to render html as string
	var walletInsideViewData views.WalletInsideViewData
	walletInsideViewData.WalletID = walletID
	walletInsideViewData.ActiveUrl = action

	t := bytes.NewBuffer([]byte{})

	if action == "/" || action == "/positions" {

		// create an http.ResponseWriter
		views.WalletInsideView(walletInsideViewData).Render(context.TODO(), t)
	} else if action == "transactions" {
		views.WalletInsideViewTransactions(walletInsideViewData).Render(context.TODO(), t)
	} else {
		views.WalletInsideViewStaking(walletInsideViewData).Render(context.TODO(), t)
	}
	// convert parseHtml to string
	htmlString := t.String()
	return htmlString
}

func SseWalletHandler(DB *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		walletID := c.Param("walletID")
		action := c.Param("action")
		fmt.Println("action", action)
		currencyToUse := "usd"
		// get session
		session := sessions.Default(c)
		// get user user setting currency
		userCurrency := session.Get(globals.UserSettingCurrency)
		if userCurrency != nil {
			currencyToUse = userCurrency.(string)
		}
		fmt.Println("currencyToUse", currencyToUse)
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
					htmlString := SseWalletContents(DB, currencyToUse, walletID, action, action)
					streamChan <- htmlString
				case <-tickerOnceOnly.C:
					htmlString := SseWalletContents(DB, currencyToUse, walletID, action, action)
					streamChan <- htmlString
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
