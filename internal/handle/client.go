package handle

import (
	"encoding/json"
	"fmt"
	global "github.com/IUnlimit/perpetua/internal"
	"github.com/IUnlimit/perpetua/internal/utils"
	"github.com/IUnlimit/perpetua/pkg/deepcopy"
	"github.com/bytedance/gopkg/util/gopool"
	log "github.com/sirupsen/logrus"
	"time"
)

type Client interface {
	writeFunc(data global.MsgData) error
	readFunc() ([]byte, error)
	getHandler() *Handler
	getUrl() string
}

// ConfigureRWWClientHandler configure and link handler to connection
func ConfigureRWWClientHandler(c *ReadWithWriteClient) {
	configureClientHandlerFunc(c, func() {
		gopool.Go(func() {
			c.handler.AddWait()
			writeAndReadClientLoop(c)
		})
	})
}

// ConfigureRAWClientHandler configure and link handler to connection
func ConfigureRAWClientHandler(c *ReadThenWriteClient) {
	configureClientHandlerFunc(c, func() {
		// write to client
		// perp -> client
		gopool.Go(func() {
			c.handler.AddWait()
			write2ClientLoop(c)
		})

		// read from client
		// perp <- client
		gopool.Go(func() {
			c.handler.AddWait()
			readFromClientLoop(c)
		})
	})
}

func configureClientHandlerFunc(c Client, rwFunc func()) {
	handler := c.getHandler()
	url := c.getUrl()
	handleSet.Add(handler)
	handler.AddWait()
	err := c.writeFunc(global.Lifecycle)
	if err != nil {
		log.Errorf("Error occurred when send lifecycle event: %v", err)
		handler.WaitExitAll()
		return
	}
	log.Infof("[Client] Connection established with url: %s", url)
	ClientOnlineStatusChangeEvent(handler, true)

	// heartbeat
	gopool.Go(func() {
		handler.AddWait()
		err := doHeartbeat(c)
		if err != nil {
			log.Infof("[Client] Error occurred when heartbeat on %s: %v", url, err)
			handler.WaitExitAll()
		}
	})

	// read & write function invoke
	rwFunc()

	handler.WaitDone()
	log.Infof("[Client] Connection closed with addr: %s", url)
}

// loop for read and write to client
func writeAndReadClientLoop(c *ReadWithWriteClient) {
	handler := c.getHandler()
	for {
		if handler.ShouldExit() {
			log.Debugf("[<->Client] Connection write&read goroutine with addr: %s has exited", c.getUrl())
			return
		}

		handler.GetMessage(func(data global.MsgData) {
			if handler.ShouldExit() {
				return
			}
			log.Debugf("[<->Client] Try to send message to client(id-%s, name-%s)", handler.id, handler.name)
			handler.Lock.Lock()
			message, err := c.writeWithReadFunc(data)
			handler.Lock.Unlock()
			addEchoThenServe("[<->Client]", c, func() ([]byte, error) {
				return message, err
			})
		})
	}
}

// loop for write to client
func write2ClientLoop(c Client) {
	handler := c.getHandler()
	for {
		if handler.ShouldExit() {
			log.Debugf("[->Client] Connection write goroutine with addr: %s has exited", c.getUrl())
			return
		}

		handler.GetMessage(func(data global.MsgData) {
			if handler.ShouldExit() {
				return
			}
			log.Debugf("[->Client] Try to send message to client(id-%s, name-%s)", handler.id, handler.name)
			handler.Lock.Lock()
			err := c.writeFunc(data)
			handler.Lock.Unlock()
			if err != nil {
				log.Warnf("[->Client] Failed to write message to connection(url: %s): %v", c.getUrl(), err)
				handler.WaitExitAll()
			}
		})
	}
}

// loop for read from client
func readFromClientLoop(c Client) {
	handler := c.getHandler()
	for {
		if handler.ShouldExit() {
			log.Debugf("[<-Client] Connection read goroutine with url: %s has exited", c.getUrl())
			return
		}
		addEchoThenServe("[<-Client]", c, func() ([]byte, error) {
			return c.readFunc()
		})
	}
}

// try to hook request and return whether to continue loop
func interceptHookedRequests(msgData global.MsgData, c Client) (bool, error) {
	handler := c.getHandler()
	resp, exist, err := TryTouchEnhanceHook(msgData, handler)
	if err != nil {
		return true, err
	}
	// todo jump able
	if !exist {
		return false, nil
	}

	uuid, err := globalCache.Append(resp)
	if err != nil {
		handler.WaitExitAll()
		return true, err
	}

	handler.AddMessage(uuid)
	return true, nil
}

// do heartbeat to client
func doHeartbeat(c Client) error {
	handler := c.getHandler()
	impl, err := utils.GetForwardImpl()
	if err != nil {
		return err
	}
	ticker := time.NewTicker(time.Millisecond * time.Duration(impl.HeartBeatInterval))
	defer ticker.Stop()

	for {
		if handler.ShouldExit() {
			return nil
		}

		// update heartbeat time
		heartbeat := deepcopy.Copy(global.Heartbeat).(global.MsgData)
		heartbeat["time"] = time.Now().UnixMilli()
		handler.Lock.Lock()
		err := c.writeFunc(heartbeat)
		handler.Lock.Unlock()
		if err != nil {
			return err
		}

		select {
		case <-ticker.C:
			log.Debugf("Heartbeat to client: %s", c.getUrl())
		}
	}
}

func addEchoThenServe(prefix string, c Client, messageSupplier func() ([]byte, error)) {
	handler := c.getHandler()
	message, err := messageSupplier()
	if err != nil {
		log.Warnf("%s Failed to read message from connection(url: %s): %v", prefix, c.getUrl(), err)
		handler.WaitExitAll()
		return
	}
	log.Debugf("%s Received message with url-%s: %s", prefix, c.getUrl(), string(message))

	if len(message) == 0 {
		return
	}
	var msgData global.MsgData
	err = json.Unmarshal(message, &msgData)
	if err != nil {
		log.Errorf("%s Failed to unmarshal client message: %s", prefix, string(message))
		return
	}

	if msgData["echo"] == nil {
		msgData["echo"] = ""
	}

	// intercept hooked requests
	exist, err := interceptHookedRequests(msgData, c)
	if exist {
		// touch hook method error, no need to break
		if err != nil {
			log.Errorf("%s Failed to intercept hooked requests: %v", prefix, msgData)
		}
		return
	}

	// sign with echo field
	id := handler.GetId()
	echo := msgData["echo"].(string)
	echo = fmt.Sprintf("%s#%s#%s", global.EchoPrefix, id, echo)
	msgData["echo"] = echo
	log.Debugf("%s Update client(url-%s) message echo: %s", prefix, c.getUrl(), echo)
	echoMap.JustPut(id, msgData)
	echoMap.Receive <- true
}
