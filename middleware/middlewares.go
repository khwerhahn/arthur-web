package middleware

import (
	"arthur-web/auth"
	"arthur-web/globals"
	"fmt"
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/golang-module/carbon/v2"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		user, password, _ := c.Request.BasicAuth()
		// get cookie details
		session := sessions.Default(c)
		isAuthenticated := session.Get(globals.IsAuthenticated)
		isAdmin := session.Get(globals.IsAdmin)
		validUntil := session.Get(globals.ValidUntil)
		// validate
		if isAuthenticated == nil || isAdmin == nil || validUntil == nil {
			fmt.Println("not isAuthenticated or isAdmin or validUntil")
			// forward to login
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		if !isAuthenticated.(bool) {
			fmt.Println("not authenticated")
			// forward to login
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		// parse validUntil as unix
		validUntilCarbon := carbon.CreateFromTimestamp(validUntil.(int64))
		fmt.Println("validUntilCarbon: ", validUntilCarbon)
		if carbon.Now().Gt(validUntilCarbon) {
			fmt.Println("expired")
			// forward to login
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		// TODO: validate user's session or token
		isAuthed, _, err := auth.AuthenticateUser(user, password)
		if err != nil {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		if isAuthed {
			c.Next()
		} else {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}
	}
}
