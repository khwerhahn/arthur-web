package handlers

import (
	"arthur-web/views"
	"net/http"

	"github.com/gin-gonic/gin"
)

func DashboardHandler() gin.HandlerFunc {
	// create a new ViewObj

	return func(c *gin.Context) {
		dashboardViewObj := views.NewViewObj("Dashboard", "/dashboard")
		dashboardViewObj, err := dashboardViewObj.UpdateViewObjSession(c)
		if err != nil {
			c.HTML(http.StatusBadRequest, "", views.Login(dashboardViewObj))
			return
		}
		c.HTML(http.StatusOK, "", views.DashboardPage(dashboardViewObj))
		return
	}
}
