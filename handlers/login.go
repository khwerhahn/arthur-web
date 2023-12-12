package handlers

import (
	"arthur-web/auth"
	"arthur-web/globals"
	"arthur-web/views"
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/golang-module/carbon/v2"
)

func LoginHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.HTML(http.StatusOK, "", views.Login(views.NewViewObj("Login")))
	}
}

func LoginPostHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		// get user and password from form
		user := c.PostForm("user")
		password := c.PostForm("password")
		// validate user and password
		// if empty return error
		if user == "" || password == "" {
			newViewObj := views.NewViewObj("Login")
			newViewObj.AddError("form", "User or password cannot be empty")
			newViewObj.AddError("user", "User can't be empty")
			newViewObj.AddError("password", "Password can't be empty")
			c.HTML(http.StatusBadRequest, "", views.Login(newViewObj))
			return
		}

		// authenticate user
		isAuthenticated, isAdmin, err := auth.AuthenticateUser(user, password)
		if err != nil {
			newViewObj := views.NewViewObj("Login")
			newViewObj.AddError("form", "User or password incorrect")
			newViewObj.AddError("user", "User or password incorrect")
			newViewObj.AddError("password", "User or password incorrect")
			c.HTML(http.StatusBadRequest, "", views.Login(newViewObj))
			return
		} else if isAuthenticated {
			// create session and redirect to dashboard
			session := sessions.Default(c)
			session.Set(globals.Userkey, user)
			session.Set(globals.IsAuthenticated, true)
			session.Set(globals.IsAdmin, isAdmin)
			// cookie expires after 1 minute
			timeNow := carbon.Now()
			timeNow.AddMinute()
			// to unix
			timeNowString := timeNow.Timestamp()
			session.Set(globals.ValidUntil, timeNowString)
			session.Save()
			c.Redirect(http.StatusFound, "/dashboard")
			return
		} else {
			newViewObj := views.NewViewObj("Login")
			// unknown error
			newViewObj.AddError("form", "Something went wrong")
			c.HTML(http.StatusBadRequest, "", views.Login(newViewObj))
			return
		}
	}
}
