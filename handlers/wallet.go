package handlers

import (
	"arthur-web/views"
	"strconv"
	"time"

	"net/http"
	"net/url"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func WalletHandler(DB *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		// get /wallet/:walletID
		walletID := c.Param("walletID")
		action := c.Param("action")
		params := c.Request.URL.Query()

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

		paramString := ""
		// cylce through params and add to paramStrings. Add "&" if not last item in url.Values
		cycledItems := 0
		for key, value := range params {
			if cycledItems == 0 {
				paramString += key + "=" + value[0]
			} else {
				paramString += "&" + key + "=" + value[0]
			}
			cycledItems++
		}

		viewObj := views.NewViewObj("Wallet", "/wallet/"+walletID, views.Style{}, views.HTMXsse{
			Url:  "/sse/wallet/" + walletID + action + "?" + paramString,
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
