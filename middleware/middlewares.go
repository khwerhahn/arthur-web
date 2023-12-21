package middleware

import (
	"arthur-web/globals"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

func AuthMiddleware(sse bool) gin.HandlerFunc {
	return func(c *gin.Context) {
		redirectLogin := url.URL{Path: "/login"}
		// get cookie details
		session := sessions.Default(c)
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
			c.Redirect(http.StatusFound, redirectLogin.RequestURI())
			return
		}
		// parse validUntil as unix
		validUntilParsed := validUntil.(int)
		// check if cookie is expired
		unixNow := time.Now().Unix()
		if unixNow > int64(validUntilParsed) {
			fmt.Println("cookie expired")
			// forward to login
			c.Redirect(http.StatusFound, redirectLogin.RequestURI())
			return
		}

		// if sse
		if sse {
			// set response headers
			c.Writer.Header().Set("Content-Type", "text/event-stream")
			c.Writer.Header().Set("Cache-Control", "no-cache")
			c.Writer.Header().Set("Connection", "keep-alive")
			c.Writer.Header().Set("Transfer-Encoding", "chunked")
		}

		// else continue
		c.Next()

	}
}
