package handlers

import (
	"arthur-web/views"
	"fmt"
	"strconv"
	"strings"
	"time"

	"net/http"
	"net/url"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func ParamSanitizer(param string) string {
	// remove from string
	unwantedStrings := []string{"?"}
	for _, v := range unwantedStrings {
		param = strings.ReplaceAll(param, v, "")
	}
	return param
}

func WalletHandler(DB *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		// get /wallet/:walletID
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
