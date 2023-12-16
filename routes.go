package main

import (
	"arthur-web/handlers"
	"arthur-web/middleware"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// manage gin routes

func PublicRoutes(g *gin.RouterGroup) {

	g.GET("/login", handlers.LoginHandler())
	g.POST("/login", handlers.LoginPostHandler())
	g.GET("/logout", handlers.LogoutHandler())
	g.GET("/", handlers.IndexHandler())

}

func PrivateRoutes(g *gin.RouterGroup, DB *gorm.DB) {
	// Dashboard
	g.GET("/dashboard", middleware.AuthMiddleware(false), handlers.DashboardHandler())
	g.GET("/sse/navbar", middleware.AuthMiddleware(true), handlers.SseNavbar(DB))
}
