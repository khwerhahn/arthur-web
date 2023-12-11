package main

import (
	"arthur-web/handlers"

	"github.com/gin-gonic/gin"
)

// manage gin routes

func PublicRoutes(g *gin.RouterGroup) {

	g.GET("/login", handlers.LoginHandler())
	// g.POST("/login", handlers.LoginPostHandler())
	g.GET("/", handlers.IndexHandler())

}

func PrivateRoutes(g *gin.RouterGroup) {

	// Dashboard
	g.GET("/dashboard", handlers.DashboardHandler())
}
