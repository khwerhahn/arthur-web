package handlers

import (
	"arthur-web/globals"
	"arthur-web/helper"
	"arthur-web/model"
	"arthur-web/views"
	"net/http"
	"strconv"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func WalletsHandler(DB *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var walletsData views.WalletsData
		indexViewObj := views.NewViewObj("Wallets", "/wallets", views.Style{}, views.HTMXsse{})
		indexViewObj, err := indexViewObj.UpdateViewObjSession(c)
		if err != nil {
			indexViewObj.AddError("top", "Something went wrong")
			c.HTML(http.StatusBadRequest, "", views.WalletsView(indexViewObj, walletsData))
			return
		}

		session := sessions.Default(c)
		userID := session.Get(globals.UserID)
		userCurrency := session.Get(globals.UserSettingCurrency)
		var userModel model.User
		user, err := userModel.GetUserByID(DB, userID.(uint))
		if err != nil {
			indexViewObj.AddError("top", "Something went wrong")
			c.HTML(http.StatusBadRequest, "", views.WalletsView(indexViewObj, walletsData))
			return
		}
		// get accounts from user_accounts
		var userAccountsModel model.UserAccounts
		userAccounts, err := userAccountsModel.GetUserAccounts(DB, user.ID)
		if err != nil {
			indexViewObj.AddError("top", "Something went wrong")
			c.HTML(http.StatusBadRequest, "", views.WalletsView(indexViewObj, walletsData))
			return
		}

		var accounts []model.Account
		for _, userAccount := range userAccounts {
			var accountModel model.Account
			account, err := accountModel.GetAccountByID(DB, userAccount.AccountID)
			if err != nil {
				indexViewObj.AddError("top", "Something went wrong")
				c.HTML(http.StatusBadRequest, "", views.WalletsView(indexViewObj, walletsData))
				return
			}
			accounts = append(accounts, account)
		}

		// create view.WalletsWallet from each account
		var wallets []views.WalletsWallet
		for _, account := range accounts {
			var wallet views.WalletsWallet
			// ada amount
			var adaAmount int64
			var fiatAmount float64
			// convert to strings
			adaAmountString := strconv.FormatInt(adaAmount, 10)
			fiatAmountString := strconv.FormatFloat(fiatAmount, 'f', 2, 64)
			wallet.Title = account.Title
			wallet.ID = account.StakeKey
			wallet.AdaAmount = adaAmountString
			wallet.FiatAmount = fiatAmountString
			wallet.UserCurrency = helper.GetSymbol(userCurrency.(string))
			wallets = append(wallets, wallet)
		}

		walletsData.Wallets = wallets

		// get

		c.HTML(http.StatusOK, "", views.WalletsView(indexViewObj, walletsData))
		return
	}
}
