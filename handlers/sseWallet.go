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

type QueryGetStakingRewardsYear struct {
	RewardsAmount int64
	EpochEnd      time.Time
	PriceEuro     float64
	PriceUsd      float64
	ValueEuro     float64
	ValueUsd      float64
}

func GetWalletInsideViewStakingDataOverviewData(walletID string) (views.WalletInsideViewStakingOverviewData, error) {
	toReturn := views.WalletInsideViewStakingOverviewData{}
	// DB := database.DB

	// get years the wallet has staked
	// for each year get the rewards amount, epoch end, price eur, value eur
	// raw query (select ah.rewards_amount, e.epoch_end, ( select m.close from marketdata m where m.range_from <= e.epoch_end and m.range_to >= e.epoch_end) as price_eur, (( select m.close from marketdata m where m.range_from <= e.epoch_end and m.range_to >= e.epoch_end) * ah.rewards_amount) as value_eur from accounthistory ah, accounts a, epochs e where a.id = ah.account_id and e.id = ah.epoch_id_available and a.id = 1 and to_char(e.epoch_end::date, 'yyyy') = '2022')

	return toReturn, nil

}

func GetWalletInsideViewStakingData(walletID string, selectedYear string) (views.WalletInsideViewStakingData, error) {
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
	toReturn.SelectedYear = selectedYear
	if selectedYear == "all" {
		dataOverview, err := GetWalletInsideViewStakingDataOverviewData(walletID)
		if err != nil {
			fmt.Println("GetWalletInsideViewStakingData error", err)
			return toReturn, err
		}
		toReturn.AllData = dataOverview
	}

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
		stakingData, err := GetWalletInsideViewStakingData(walletID, params.Get("year"))
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

		// sanitize each param value
		for k, v := range params {
			params.Set(k, ParamSanitizer(v[0]))
		}

		if action == "/staking" {
			// if params not empty then redirect to ?year=all
			// param year must either be equal to "all" or a year between 2017 and current year
			// if year is a number then it must be between 2017 and current year
			// otherwise redirect to /wallet/walletID/staking?year=all
			if len(params) > 0 {
				// check if param "year" exists
				_, ok := params["year"]
				if ok {
					// make sure year contains string of either "all" or a number between 2017 and current year
					year := params.Get("year")
					if year != "all" {
						fmt.Println("not all")
						// check if it contains a year between 2017 and current years
						yearInt, err := strconv.Atoi(year)
						if err != nil {
							// redirect
							redirectLocation := url.URL{Path: "/wallet/" + walletID + "/staking?year=all"}
							c.Redirect(http.StatusFound, redirectLocation.RequestURI())
							return
						}
						// check if year is between 2017 and current year
						currentYear := time.Now().Year()
						if yearInt < 2017 || yearInt > currentYear {
							// redirect
							redirectLocation := url.URL{Path: "/wallet/" + walletID + "/staking?year=all"}
							c.Redirect(http.StatusFound, redirectLocation.RequestURI())
							return
						}
					}
				} else {
					fmt.Println("no year")
					// redirect
					redirectLocation := url.URL{Path: "/wallet/" + walletID + "/staking?year=all"}
					c.Redirect(http.StatusFound, redirectLocation.RequestURI())
					return
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
