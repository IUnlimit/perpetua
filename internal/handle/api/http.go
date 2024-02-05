package api

import (
	"fmt"
	"github.com/IUnlimit/perpetua/internal/handle"
	"github.com/IUnlimit/perpetua/internal/logger"
	"github.com/IUnlimit/perpetua/internal/utils"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"io"
	"net"
	"os"
)

// Custom API impl

// GetWebSocketPort get available ws port
func GetWebSocketPort(ctx *gin.Context) {
	listen, err := net.Listen("tcp", ":0")
	if err != nil {
		utils.BadResponse(ctx, err)
		return
	}
	_ = listen.Close()

	port := listen.Addr().(*net.TCPAddr).Port
	// open new ws service, wait for minutes
	handle.CreateWSInstance(port)
	log.Info("Create new websocket connection on port: ", port)

	utils.SendResponse(ctx,
		"port", port)
}

// EnableHttpService enable http server
func EnableHttpService(port int) {
	gin.SetMode(gin.ReleaseMode)
	gin.DefaultWriter = io.MultiWriter(os.Stdout, logger.Hook.GetWriter())
	engine := gin.Default()

	router := engine.Group("/api")
	router.GET("/get_ws_port", GetWebSocketPort)

	log.Info("Starting the HTTP service on port: ", port)
	err := engine.Run(fmt.Sprintf(":%d", port))
	if err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
}
