package handle

import (
	"context"
	"fmt"
	"github.com/IUnlimit/perpetua/internal/conf"
	"github.com/gorilla/websocket"
	log "github.com/sirupsen/logrus"
	"net/http"
	"time"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		// allow all sources conn
		return true
	},
}

func CreateWSInstance(port int) {
	var start bool
	handleWebSocket := func(w http.ResponseWriter, r *http.Request) {
		// upgrade HTTP to WebSocket conn
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Errorf("Failed to upgrade connection to WebSocket: %v", err)
			return
		}
		defer conn.Close()

		start = true
		log.Info("WebSocket connection established on port: ", port)
		for {
			messageType, message, err := conn.ReadMessage()
			if err != nil {
				log.Warnf("Failed to read message from WebSocket(port: %d): %v", port, err)
				break
			}

			log.Debugf("Received message(type: %d) on port-%d: %s", messageType, port, string(message))
			// TODO redirect
			var resp []byte

			err = conn.WriteMessage(messageType, resp)
			if err != nil {
				log.Warnf("Failed to write message to WebSocket(port: %d): %v", port, err)
				break
			}
		}
	}

	server := &http.Server{
		Addr:    fmt.Sprintf(":%d", port),
		Handler: http.HandlerFunc(handleWebSocket),
	}

	config := conf.Config.WebSocket
	go func() {
		timer := time.After(config.Timeout)
		<-timer
		if !start {
			_ = server.Shutdown(context.Background())
		}
	}()

	err := server.ListenAndServe()
	if err != nil {
		log.Warnf("WebSocket connection(port: %d) exit with error: %v", port, err)
	}
}

// 处理 WebSocket 请求
//func handleWebSocket(w http.ResponseWriter, r *http.Request) {
//	// 建立 WebSocket 连接
//	upgrader := websocket.Upgrader{}
//	conn, err := upgrader.Upgrade(w, r, nil)
//	if err != nil {
//		log.Println("Failed to upgrade to WebSocket:", err)
//		return
//	}
//	defer conn.Close()
//
//	// 连接到目标 WebSocket 服务器
//	targetURL := "ws://localhost:8081/ws"
//	targetConn, _, err := websocket.DefaultDialer.Dial(targetURL, nil)
//	if err != nil {
//		log.Println("Failed to connect to target WebSocket server:", err)
//		return
//	}
//	defer targetConn.Close()
//
//	// 在两个 WebSocket 连接之间进行转发
//	for {
//		// 从客户端读取消息
//		_, clientMsg, err := conn.ReadMessage()
//		if err != nil {
//			log.Println("Failed to read client message:", err)
//			break
//		}
//
//		// 转发消息到目标 WebSocket 服务器
//		err = targetConn.WriteMessage(websocket.TextMessage, clientMsg)
//		if err != nil {
//			log.Println("Failed to write message to target WebSocket server:", err)
//			break
//		}
//
//		// 从目标 WebSocket 服务器读取响应
//		_, targetMsg, err := targetConn.ReadMessage()
//		if err != nil {
//			log.Println("Failed to read target server message:", err)
//			break
//		}
//
//		// 将响应消息发送给客户端
//		err = conn.WriteMessage(websocket.TextMessage, targetMsg)
//		if err != nil {
//			log.Println("Failed to write message to client:", err)
//			break
//		}
//	}
//}
