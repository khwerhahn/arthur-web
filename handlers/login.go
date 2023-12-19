package handlers

import (
	"arthur-web/auth"
	"arthur-web/globals"
	"arthur-web/views"
	"net/http"
	"time"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

func LoginHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		loginStyle := views.Style{
			StyleContainer: []string{"items-center", "justify-center", "flex-1"},
		}
		c.HTML(http.StatusOK, "", views.Login(views.NewViewObj("Login", "/login", loginStyle)))
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
			newViewObj := views.NewViewObj("Login", "/login", views.Style{})
			newViewObj.AddError("form", "User or password cannot be empty")
			newViewObj.AddError("user", "User can't be empty")
			newViewObj.AddError("password", "Password can't be empty")
			c.HTML(http.StatusBadRequest, "", views.Login(newViewObj))
			return
		}

		// authenticate user
		userDB, err := auth.AuthenticateUser(user, password)
		if err != nil {
			newViewObj := views.NewViewObj("Login", "/login", views.Style{})
			newViewObj.AddError("form", "User or password incorrect")
			newViewObj.AddError("user", "User or password incorrect")
			newViewObj.AddError("password", "User or password incorrect")
			c.HTML(http.StatusBadRequest, "", views.Login(newViewObj))
			return
		} else if userDB.ID != 0 {
			// create session and redirect to dashboard
			session := sessions.Default(c)
			session.Set(globals.Userkey, user)
			session.Set(globals.IsAuthenticated, true)
			session.Set(globals.IsAdmin, userDB.IsAdmin)
			session.Set(globals.UserID, userDB.ID)
			session.Set(globals.ProfileImageUrl, userDB.ProfileImageUrl)
			// extract user settings
			userSettings, err := userDB.GetUserSettings()
			if err != nil {
				newViewObj := views.NewViewObj("Login", "/login", views.Style{})
				newViewObj.AddError("form", "Something went wrong")
				c.HTML(http.StatusBadRequest, "", views.Login(newViewObj))
				return

			}
			session.Set(globals.UserSettingCurrency, userSettings.Currency)
			session.Set(globals.UserSettingLanguage, userSettings.Language)
			// cookie expires in minutes between now and midnight
			diffNowToMidnight := time.Until(time.Now().AddDate(0, 0, 1).Truncate(24 * time.Hour))
			expirtationTime := time.Now().Add(diffNowToMidnight * time.Minute)
			session.Set(globals.ValidUntil, int(expirtationTime.Unix()))
			session.Save()
			c.Redirect(http.StatusFound, "/dashboard")
			return
		} else {
			newViewObj := views.NewViewObj("Login", "/login", views.Style{})
			// unknown error
			newViewObj.AddError("form", "Something went wrong")
			c.HTML(http.StatusBadRequest, "", views.Login(newViewObj))
			return
		}
	}
}
