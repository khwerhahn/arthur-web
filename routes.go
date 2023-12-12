package main

import (
	"arthur-web/handlers"
	"arthur-web/middleware"

	"github.com/gin-gonic/gin"
)

// manage gin routes

func PublicRoutes(g *gin.RouterGroup) {

	g.GET("/login", handlers.LoginHandler())
	g.POST("/login", handlers.LoginPostHandler())
	g.GET("/", handlers.IndexHandler())

}

func PrivateRoutes(g *gin.RouterGroup) {
	// all routes in this group are private and require AuthMiddleware
	g.Use(middleware.AuthMiddleware())
	// Dashboard
	g.GET("/dashboard", handlers.DashboardHandler())
}
