package web

import (
	"fmt"
	"io"
	"io/fs"
	"net/http"
	"os"
	"strings"

	global "github.com/IUnlimit/perpetua/internal"
	"github.com/IUnlimit/perpetua/internal/logger"
	webstatic "github.com/IUnlimit/perpetua/web"
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
	api.GET("/packets/trace", GetPacketTrace)
	api.GET("/system", GetSystemInfo)
	api.DELETE("/packets", DeletePackets)

	// Mount embedded assets subdirectory at /assets
	assetsFS, err := fs.Sub(webstatic.StaticFS, "assets")
	if err != nil {
		log.Errorf("[Web] Failed to load embedded assets: %v", err)
		return
	}
	engine.StaticFS("/assets", http.FS(assetsFS))

	// Read index.html once from embed
	indexHTML, err := webstatic.StaticFS.ReadFile("index.html")
	if err != nil {
		log.Errorf("[Web] Failed to read embedded index.html: %v", err)
		return
	}

	serveIndex := func(c *gin.Context) {
		c.Data(http.StatusOK, "text/html; charset=utf-8", indexHTML)
	}

	// Root route
	engine.GET("/", serveIndex)

	// Fallback: serve index.html for non-API, non-asset routes
	engine.NoRoute(func(c *gin.Context) {
		path := c.Request.URL.Path
		if strings.HasPrefix(path, "/api/") || strings.HasPrefix(path, "/assets/") {
			c.Status(http.StatusNotFound)
			return
		}
		serveIndex(c)
	})

	log.Infof("[Web] Starting management dashboard on port: %d", port)
	go func() {
		if err := engine.Run(fmt.Sprintf(":%d", port)); err != nil && err != http.ErrServerClosed {
			log.Errorf("[Web] Failed to start web server: %v", err)
		}
	}()
}
