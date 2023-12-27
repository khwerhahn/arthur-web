package handlers

import (
	"arthur-web/database"
	"arthur-web/globals"
	"arthur-web/model"
	"arthur-web/views"
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func GetWalletInsideViewStakingData(walletID string) (views.WalletInsideViewStakingData, error) {
	DB := database.DB
	toReturn := views.WalletInsideViewStakingData{}

	// get available years
	var availableAccountYearsQuery []views.AccountAvailableYears
	// raw Query
	err := DB.Raw("select distinct to_char(e.epoch_end, 'YYYY') as year from epochs e, accounthistory h, accounts a where e.id = h.epoch_id and a.stake_key = ? and h.account_id = a.id order by year asc", walletID).Scan(&availableAccountYearsQuery).Error
	if err != nil {
		fmt.Println("GetWalletInsideViewStakingData error", err)
		return toReturn, err
	}
	// add to toReturn
	toReturn.AvailableYears = availableAccountYearsQuery

	return toReturn, nil
}

func SseWalletContents(DB *gorm.DB, currencyToUse string, walletID string, activeUrl string, action string, params url.Values) string {

	// get account (wallet) details
	var accountData model.Account
	err := accountData.GetAccountByStakeKey(DB, walletID)
	if err != nil {
		fmt.Println("SseWalletContents error", err)
		return "error"
	}

	// use templ to render html as string
	var walletInsideViewData views.WalletInsideViewData
	walletInsideViewData.WalletID = walletID
	walletInsideViewData.WalletTitle = accountData.Title
	walletInsideViewData.ActiveUrl = action

	t := bytes.NewBuffer([]byte{})

	if action == "/" || action == "/positions" {

		// create an http.ResponseWriter
		views.WalletInsideView(walletInsideViewData).Render(context.TODO(), t)
	} else if action == "transactions" {
		views.WalletInsideViewTransactions(walletInsideViewData).Render(context.TODO(), t)
	} else {
		// Staking

		stakingData, err := GetWalletInsideViewStakingData(walletID)
		if err != nil {
			// write eroor to t
			t.WriteString("error")
		}

		views.WalletInsideViewStaking(walletInsideViewData, stakingData).Render(context.TODO(), t)
	}
	// convert parseHtml to string
	htmlString := t.String()
	return htmlString
}

func SseWalletHandler(DB *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		walletID := c.Param("walletID")
		action := c.Param("action")
		params := c.Request.URL.Query()
		// allowed params values are "all" and year numbers starting frokm 2017 to current year

		if action == "/staking" {

			// if params not empty or contains "all" or year then redirect to ?year=all
			// if year is a number then it must be between 2017 and current year
			if len(params) > 0 {
				// check if params contains "all" or year
				_, ok := params["all"]
				if ok {
					// redirect to ?year=all
					redirectLocation := url.URL{Path: "/wallet/" + walletID + "/staking?year=all"}
					c.Redirect(http.StatusFound, redirectLocation.RequestURI())
					return
				}
				// check if params contains year
				year, ok := params["year"]
				if ok {
					// check if year is a number
					yearInt, err := strconv.Atoi(year[0])
					if err == nil {
						// check if year is between 2017 and current year
						redirectLocation := url.URL{Path: "/wallet/" + walletID + "/staking?year=all"}
						c.Redirect(http.StatusFound, redirectLocation.RequestURI())
						return
					}
					if yearInt >= 2017 && yearInt <= time.Now().Year() {
						redirectLocation := url.URL{Path: "/wallet/" + walletID + "/staking?year=all"}
						c.Redirect(http.StatusFound, redirectLocation.RequestURI())
						return
					}
				}
			} else {
				// redirect to ?year=all
				redirectLocation := url.URL{Path: "/wallet/" + walletID + "/staking?year=all"}
				c.Redirect(http.StatusFound, redirectLocation.RequestURI())
				return
			}
		}

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
					htmlString := SseWalletContents(DB, currencyToUse, walletID, action, action, params)
					streamChan <- htmlString
				case <-tickerOnceOnly.C:
					htmlString := SseWalletContents(DB, currencyToUse, walletID, action, action, params)
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
