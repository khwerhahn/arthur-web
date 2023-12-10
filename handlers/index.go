package handlers

// handerls (controllers) for gin routes

import (
	"arthur-web/views"
	"net/http"

	"github.com/gin-gonic/gin"
)

// IndexHandler handles the index routes
func IndexHandler() gin.HandlerFunc {
	// create a new ViewObj
	indexViewObj := views.NewViewObj("Index")
	return func(c *gin.Context) {
		c.HTML(http.StatusOK, "", views.IndexPage(indexViewObj))
	}
}
