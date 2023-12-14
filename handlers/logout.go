package handlers

import (
	"arthur-web/views"
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

func LogoutHandler() gin.HandlerFunc {
	// create a new ViewObj
	return func(c *gin.Context) {
		// destroy session
		session := sessions.Default(c)
		session.Clear()
		session.Save()

		logoutViewObj := views.NewViewObj("Logout", "/logout")
		logoutViewObj, err := logoutViewObj.UpdateViewObjSession(c)
		if err != nil {
			c.HTML(http.StatusBadRequest, "", views.Logout(logoutViewObj))
			return
		}
		c.HTML(http.StatusOK, "", views.Logout(logoutViewObj))
		return
	}
}
