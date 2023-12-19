package handlers

import (
	"arthur-web/database"
	"arthur-web/globals"
	"arthur-web/model"
	"arthur-web/views"
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

func DashboardHandler() gin.HandlerFunc {
	// create a new ViewObj

	return func(c *gin.Context) {
		DB := database.DB

		var dashboardData views.DashboardData

		dashboardViewObj := views.NewViewObj("Dashboard", "/dashboard", views.Style{})
		dashboardViewObj, err := dashboardViewObj.UpdateViewObjSession(c)

		// get accounts from user_accounts
		// get user from UpdateViewObjSession
		session := sessions.Default(c)
		userID := session.Get(globals.UserID)
		var userModel model.User
		user, err := userModel.GetUserByID(DB, userID.(uint))
		if err != nil {
			dashboardViewObj.AddError("top", "Something went wrong")
			c.HTML(http.StatusBadRequest, "", views.DashboardPage(dashboardViewObj, dashboardData))
			return
		}
		// get accounts from user_accounts
		var userAccountsModel model.UsersAccounts
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
			var dashboardWallet views.DashboardWallet
			dashboardWallet.ID = account.ID
			dashboardWallet.Title = account.Title
			dashboardWallet.ADAAmount = 0
			dashboardWallet.FiatAmount = 0
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
