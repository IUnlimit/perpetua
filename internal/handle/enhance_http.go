package handle

import (
	"fmt"
	global "github.com/IUnlimit/perpetua/internal"
	"github.com/IUnlimit/perpetua/internal/logger"
	"github.com/IUnlimit/perpetua/internal/model"
	"github.com/IUnlimit/perpetua/internal/utils"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"io"
	"os"
)

// Custom API impl

// GetWebSocketPort get available ws port
func GetWebSocketPort(ctx *gin.Context) {
	rangePort := global.Config.WebSocket.RangePort
	port, err := utils.RandomAvailablePort(rangePort.Enabled, rangePort.Start, rangePort.End)
	if err != nil {
		utils.BadResponse(ctx, err)
		return
	}

	// open new ws service, wait for minutes
	CreateWSInstance(port)
	log.Info("[Enhance] Create new websocket connection on port: ", port)

	utils.GoodResponse(ctx,
		"port", port)
}

// GetOnlineClients get online clients array
func GetOnlineClients(ctx *gin.Context) {
	clients := make([]*model.Client, 0)
	for _, v := range handleSet.Iterator() {
		handler := v.(*Handler)
		clients = append(clients, &model.Client{
			AppId:      handler.GetId(),
			ClientName: handler.GetName(),
		})
	}
	utils.GoodResponseArray(ctx,
		clients)
}

// EnableHttpService enable http server
func EnableHttpService(port int) {
	gin.SetMode(gin.ReleaseMode)
	gin.DefaultWriter = io.MultiWriter(os.Stdout, logger.Hook.GetWriter())
	engine := gin.Default()

	router := engine.Group("/api")
	router.GET("/get_ws_port", GetWebSocketPort)
	router.GET("/get_online_clients", GetOnlineClients)

	log.Info("[Enhance] Starting the HTTP service on port: ", port)
	err := engine.Run(fmt.Sprintf(":%d", port))
	if err != nil {
		log.Fatalf("[Enhance] failed to start server: %v", err)
	}
}
