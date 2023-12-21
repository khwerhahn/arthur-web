package handlers

import (
	"arthur-web/views"
	"fmt"

	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func WalletHandler(DB *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		// get /wallet/:walletID
		walletID := c.Param("walletID")
		action := c.Param("action")
		fmt.Println("WalletHandler", walletID, action)


		viewObj := views.NewViewObj("Wallet", "/wallet/"+walletID, views.Style{}, views.HTMXsse{
			Url:  "/sse/wallet/" + walletID + action,
			Swap: "message",
		})
		viewObj, err := viewObj.UpdateViewObjSession(c)
		// WalletViewData
		var walletViewData views.WalletViewData
		walletViewData.WalletID = walletID
		if err != nil {
			c.HTML(http.StatusBadRequest, "", views.WalletView(viewObj, walletViewData))
			return
		}
		c.HTML(http.StatusOK, "", views.WalletView(viewObj, walletViewData))
		return
	}
}


