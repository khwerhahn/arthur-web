package main

import (
	"arthur-web/config"
	"arthur-web/handlers"
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
)

func main() {
	ginMode := ""
	ginModeEnv := config.Config("GIN_MODE")
	if ginModeEnv == "release" {
		ginMode = gin.ReleaseMode
	} else {
		ginMode = gin.DebugMode
	}

	gin.SetMode(ginMode)
	router := gin.Default()
	router.HTMLRender = &TemplRender{}
	var sessionSecret []byte
	secretEnv := config.Config("SECRET")
	if len(secretEnv) > 0 {
		sessionSecret = []byte(secretEnv)
	} else {
		sessionSecret = []byte("secret")
	}
	router.Use(sessions.Sessions("session", cookie.NewStore(sessionSecret)))

	//////////////////////
	// Routes
	// serve static files
	router.Static("/assets", "./assets")
	// 404 Handler
	router.NoRoute(handlers.NotFoundHandler)
	public := router.Group("/")
	PublicRoutes(public)
	private := router.Group("/")
	PrivateRoutes(private)
	//////////////////////

	srv := &http.Server{
		Addr:    ":8080",
		Handler: router,
	}

	go func() {
		// service connections
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server with
	// a timeout of 5 seconds.
	quit := make(chan os.Signal)
	// kill (no param) default send syscanll.SIGTERM
	// kill -2 is syscall.SIGINT
	// kill -9 is syscall. SIGKILL but can"t be catch, so don't need add it
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutdown Server ...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server Shutdown:", err)
	}
	// catching ctx.Done(). timeout of 5 seconds.
	select {
	case <-ctx.Done():
		log.Println("timeout of 5 seconds.")
	}
	log.Println("Server exiting")
}
