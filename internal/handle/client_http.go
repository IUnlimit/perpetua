package handle

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	global "github.com/IUnlimit/perpetua/internal"
	"github.com/IUnlimit/perpetua/internal/utils"
	log "github.com/sirupsen/logrus"
	"io"
	"net/http"
	"strconv"
	"time"
)

// ReadWithWriteClient read loop with write loop
type ReadWithWriteClient struct {
	url     string
	secret  string
	selfId  string
	handler *Handler
}

func NewReadWithWriteClient(url string, secret string, handler *Handler) *ReadWithWriteClient {
	return &ReadWithWriteClient{
		url:     url,
		secret:  secret,
		selfId:  strconv.Itoa(int(global.Lifecycle["self_id"].(float64))),
		handler: handler,
	}
}

func (rww *ReadWithWriteClient) writeWithReadFunc(data global.MsgData) ([]byte, error) {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest("POST", rww.url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Self-ID", rww.selfId)
	if len(rww.secret) != 0 {
		req.Header.Set("X-Signature", "")
	}

	resp, err := (&http.Client{}).Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if !(200 <= resp.StatusCode && resp.StatusCode <= 299) {
		return nil, errors.New(fmt.Sprintf("error response status: %d", resp.StatusCode))
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return body, nil
}

func (rww *ReadWithWriteClient) writeFunc(data global.MsgData) error {
	_, err := rww.writeWithReadFunc(data)
	return err
}

func (rww *ReadWithWriteClient) readFunc() ([]byte, error) {
	return nil, errors.New("unsupported operation")
}

func (rww *ReadWithWriteClient) getHandler() *Handler {
	return rww.handler
}

func (rww *ReadWithWriteClient) getUrl() string {
	return rww.url
}

func TryPostHttp(postUrl string, secret string) error {
	impl, err := utils.GetForwardImpl()
	if err != nil {
		return err
	}
	<-utils.WaitNTQQStartup(impl.Host, impl.Port, nil)
	<-utils.WaitCondition(time.Duration(2000), func() error {
		if global.Lifecycle == nil {
			return errors.New("not init yet")
		}
		return nil
	}, nil)

	log.Infof("[Client] Start try to report events to postUrl: %s", postUrl)
	handler := NewHandler(context.Background())
	client := NewReadWithWriteClient(postUrl, secret, handler)
	ConfigureRWWClientHandler(client)
	return nil
}
