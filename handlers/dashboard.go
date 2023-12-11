package handlers

import (
	"arthur-web/views"
	"net/http"

	"github.com/gin-gonic/gin"
)

func DashboardHandler() gin.HandlerFunc {
	// create a new ViewObj
	dashboardViewObj := views.NewViewObj("Dashboard")
	return func(c *gin.Context) {
		c.HTML(http.StatusOK, "", views.DashboardPage(dashboardViewObj))
	}
}
