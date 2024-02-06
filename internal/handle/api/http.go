package api

import (
	"errors"
	"fmt"
	global "github.com/IUnlimit/perpetua/internal"
	"github.com/IUnlimit/perpetua/internal/handle"
	"github.com/IUnlimit/perpetua/internal/logger"
	"github.com/IUnlimit/perpetua/internal/model"
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
	rangePort := global.Config.WebSocket.RangePort
	port, err := randomAvailablePort(rangePort)
	if err != nil {
		utils.BadResponse(ctx, err)
		return
	}

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

func randomAvailablePort(rangePort *model.RangePort) (int, error) {
	if !rangePort.Enabled {
		port, err := tryListen(0)
		if err != nil {
			return 0, err
		}
		return port, nil
	}

	portRange := rangePort.End - rangePort.Start
	if rangePort.Start < 512 || rangePort.End > 65535 || portRange < 0 {
		return 0, errors.New(fmt.Sprintf("invalid port range [%d-%d]", rangePort.Start, rangePort.End))
	}
	if portRange == 0 {
		return rangePort.Start, nil
	}

	for i := rangePort.Start; i <= rangePort.End; i++ {
		port, err := tryListen(i)
		if err == nil {
			return port, nil
		}
	}
	return 0, errors.New(fmt.Sprintf("unavailable port in range [%d-%d]", rangePort.Start, rangePort.End))
}

func tryListen(port int) (int, error) {
	listen, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return 0, err
	}
	_ = listen.Close()

	port = listen.Addr().(*net.TCPAddr).Port
	return port, nil
}
