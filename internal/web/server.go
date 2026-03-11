package web

import (
	"fmt"
	"io"
	"net/http"
	"os"

	global "github.com/IUnlimit/perpetua/internal"
	"github.com/IUnlimit/perpetua/internal/logger"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

// StartWebServer starts the web management dashboard server
func StartWebServer() {
	port := 9090
	if global.Config.Web != nil {
		port = global.Config.Web.Port
	}

	gin.SetMode(gin.ReleaseMode)
	gin.DefaultWriter = io.MultiWriter(os.Stdout, logger.Hook.GetWriter())
	engine := gin.Default()

	// CORS middleware
	engine.Use(func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Content-Type")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	})

	// API routes
	api := engine.Group("/api/web")
	api.GET("/connections", GetConnections)
	api.GET("/packets", GetPackets)
	api.GET("/system", GetSystemInfo)
	api.DELETE("/packets", DeletePackets)

	// Serve static frontend files
	engine.Static("/assets", "./web/assets")
	engine.StaticFile("/", "./web/index.html")
	engine.StaticFile("/favicon.ico", "./web/favicon.ico")
	// Fallback for SPA-like routing
	engine.NoRoute(func(c *gin.Context) {
		c.File("./web/index.html")
	})

	log.Infof("[Web] Starting management dashboard on port: %d", port)
	go func() {
		if err := engine.Run(fmt.Sprintf(":%d", port)); err != nil && err != http.ErrServerClosed {
			log.Errorf("[Web] Failed to start web server: %v", err)
		}
	}()
}
