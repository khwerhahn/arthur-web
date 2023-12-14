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
	return func(c *gin.Context) {
		indexViewObj := views.NewViewObj("Index", "/")
		indexViewObj, err := indexViewObj.UpdateViewObjSession(c)
		if err != nil {
			c.HTML(http.StatusBadRequest, "", views.Login(indexViewObj))
			return
		}
		c.HTML(http.StatusOK, "", views.IndexPage(indexViewObj))
		return
	}
}
