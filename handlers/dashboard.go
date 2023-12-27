package handlers

import (
	"arthur-web/database"
	"arthur-web/globals"
	"arthur-web/helper"
	"arthur-web/model"
	"arthur-web/views"
	"net/http"
	"strconv"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

func DashboardHandler() gin.HandlerFunc {
	// create a new ViewObj

	return func(c *gin.Context) {
		DB := database.DB

		var dashboardData views.DashboardData

		dashboardViewObj := views.NewViewObj("Dashboard", "/dashboard", views.Style{}, views.HTMXsse{})
		dashboardViewObj, err := dashboardViewObj.UpdateViewObjSession(c)

		// get accounts from user_accounts
		// get user from UpdateViewObjSession
		session := sessions.Default(c)
		userID := session.Get(globals.UserID)
		userCurrency := session.Get(globals.UserSettingCurrency)
		var userModel model.User
		user, err := userModel.GetUserByID(DB, userID.(uint))
		if err != nil {
			dashboardViewObj.AddError("top", "Something went wrong")
			c.HTML(http.StatusBadRequest, "", views.DashboardPage(dashboardViewObj, dashboardData))
			return
		}
		// get accounts from user_accounts
		var userAccountsModel model.UserAccounts
		userAccounts, err := userAccountsModel.GetUserAccounts(DB, user.ID)
		if err != nil {
			dashboardViewObj.AddError("top", "Something went wrong")
			c.HTML(http.StatusBadRequest, "", views.DashboardPage(dashboardViewObj, dashboardData))
			return
		}

		var accounts []model.Account
		for _, userAccount := range userAccounts {
			var accountModel model.Account
			account, err := accountModel.GetAccountByID(DB, userAccount.AccountID)
			if err != nil {
				dashboardViewObj.AddError("top", "Something went wrong")
				c.HTML(http.StatusBadRequest, "", views.DashboardPage(dashboardViewObj, dashboardData))
				return
			}
			accounts = append(accounts, account)
		}

		// create view.DashboardWallet for each account
		var dashboardWallets []views.DashboardWallet
		for _, account := range accounts {
			// ada amount
			var adaAmount int64
			var fiatAmount float64
			// convert adaAmount to string
			adaAmountString := strconv.FormatInt(adaAmount, 10)
			fiatAmountString := strconv.FormatFloat(fiatAmount, 'f', 2, 64)
			var dashboardWallet views.DashboardWallet
			dashboardWallet.ID = account.StakeKey
			dashboardWallet.Title = account.Title
			dashboardWallet.AdaAmount = adaAmountString
			dashboardWallet.FiatAmount = fiatAmountString
			dashboardWallet.UserCurrency = helper.GetSymbol(userCurrency.(string))
			dashboardWallets = append(dashboardWallets, dashboardWallet)
		}

		// append dashboardWallets to dashboardData
		dashboardData.Wallets = dashboardWallets

		if err != nil {
			c.HTML(http.StatusBadRequest, "", views.DashboardPage(dashboardViewObj, dashboardData))
			return
		}
		c.HTML(http.StatusOK, "", views.DashboardPage(dashboardViewObj, dashboardData))
		return
	}
}
