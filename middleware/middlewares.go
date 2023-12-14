package middleware

import (
	"arthur-web/globals"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// get cookie details
		session := sessions.Default(c)
		fmt.Println("session: ", session)
		isAuthenticated := session.Get(globals.IsAuthenticated)
		isAdmin := session.Get(globals.IsAdmin)
		validUntil := session.Get(globals.ValidUntil)
		// validate
		if isAuthenticated == nil || isAdmin == nil || validUntil == nil {
			fmt.Println("not isAuthenticated or isAdmin or validUntil")
			// forward to login
			c.Redirect(http.StatusFound, "/login")
			return
		}
		if !isAuthenticated.(bool) {
			fmt.Println("not authenticated")
			// forward to login
			c.Redirect(http.StatusFound, "/login")
			return
		}
		// parse validUntil as unix
		validUntilParsed := validUntil.(int)
		// check if cookie is expired
		unixNow := time.Now().Unix()
		if unixNow > int64(validUntilParsed) {
			fmt.Println("cookie expired")
			// forward to login
			c.Redirect(http.StatusFound, "/login")
			return
		}

		// else continue
		c.Next()

	}
}
