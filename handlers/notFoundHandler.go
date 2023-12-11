package handlers

import (
	"arthur-web/views"
	"net/http"

	"github.com/gin-gonic/gin"
)

// 404 handler

func NotFoundHandler(c *gin.Context) {
	viewObj := views.NewViewObj("404")
	c.HTML(http.StatusNotFound, "", views.NotFoundPage(viewObj))
}
