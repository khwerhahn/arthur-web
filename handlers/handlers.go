package handlers

import (
	"arthur-web/views"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

func GetSessionData(c *gin.Context, viewObj *views.ViewObj) (*views.ViewObj, error) {
	session := sessions.Default(c)
	isAuthenticated := session.Get("isAuthenticated")
	isAdmin := session.Get("isAdmin")

	// fill ViewObj Session
	viewObj.Session.IsAuthenticated = isAuthenticated.(bool)
	viewObj.Session.IsAdmin = isAdmin.(bool)
	viewObj.Session.Username = "NIY"
	viewObj.Session.FirstName = "NIY"
	viewObj.Session.LastName = "NIY"
	return viewObj, nil
}
