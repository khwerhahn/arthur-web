package handlers

// handerls (controllers) for gin routes

import (
	"arthur-web/views"
	"net/http"
	"net/url"

	"github.com/gin-gonic/gin"
)

// IndexHandler handles the index routes
func IndexHandler() gin.HandlerFunc {
	// create a new ViewObj
	return func(c *gin.Context) {
		indexViewObj := views.NewViewObj("Index", "/", views.Style{}, views.HTMXsse{})
		indexViewObj, err := indexViewObj.UpdateViewObjSession(c)
		indexViewObj.Style.StyleContainer = append(indexViewObj.Style.StyleContainer, "flex-1")
		if err != nil {
			indexViewObj.AddError("top", "Something went wrong")
			c.HTML(http.StatusBadRequest, "", views.IndexPage(indexViewObj))
			return
		}
		if indexViewObj.Session.IsAuthenticated {
			redirectDashboared := url.URL{Path: "/dashboard"}
			c.Redirect(http.StatusFound, redirectDashboared.RequestURI())
			return
		}
		c.HTML(http.StatusOK, "", views.IndexPage(indexViewObj))
		return
	}
}
