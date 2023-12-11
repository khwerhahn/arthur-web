package middleware

import (
	"arthur-web/auth"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func RequestLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		fmt.Println("///////////////////////////////////////////")
		fmt.Println("///////////// RequestLogger ///////////////")
		fmt.Println("Host", c.Request.Host)
		fmt.Println("RemoteAddr", c.Request.RemoteAddr)
		fmt.Println("RequestURI", c.Request.RequestURI)
		fmt.Println("Cookies", c.Request.Cookies())
		fmt.Println("Header", c.Request.Header)
		fmt.Println("///////////////////////////////////////////")
		c.Next()
	}
}

func TokenAuthMiddleware() gin.HandlerFunc {
	errList := make(map[string]string)
	return func(c *gin.Context) {
		err := auth.CheckCookie(c)
		if err != nil {
			errList["unauthorized"] = "Unauthorized"
			c.JSON(http.StatusUnauthorized, gin.H{
				"status": http.StatusUnauthorized,
				"error":  errList,
			})
			c.Abort()
			return
		}
		c.Next()
	}
}
