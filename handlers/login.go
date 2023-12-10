package handlers

import (
	"arthur-web/globals"
	"arthur-web/views"
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

func LoginHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		session := sessions.Default(c)
		user := session.Get(globals.Userkey)
		loginViewObj := views.NewViewObj("Login")
		if user != nil {
			c.HTML(http.StatusBadRequest, "", views.Login(loginViewObj))
		}
		c.HTML(http.StatusOK, "", views.Login(loginViewObj))
	}
}
