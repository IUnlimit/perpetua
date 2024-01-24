package handle

import (
	"context"
	global "github.com/IUnlimit/perpetua/internal"
	collections "github.com/chenjiandongx/go-queue"
	"sync"
)

// handleList stores handle.Handler for client websocket
var handleList []*Handler

type Handler struct {
	ctx     context.Context
	wg      sync.WaitGroup
	receive chan bool
	// waiting goroutine count
	waitCount int
	queue     *collections.Queue
}

func (h *Handler) AddMessage(uuid string) {
	h.queue.Put(uuid)
}

// GetMessage from local cache
func (h *Handler) GetMessage(invoke func(data *global.MsgData)) {
	for e, _ := h.queue.Get(); e != nil; {
		data, _ := globalCache.cache.Get(e)
		if data == nil {
			continue
		}
		invoke(data.(*global.MsgData))
	}
}

// ShouldExit 是否需要结束
func (h *Handler) ShouldExit() bool {
	return h.waitCount == 0
}

// WaitExitAll 结束所有等待
func (h *Handler) WaitExitAll() {
	for i := 0; i < h.waitCount; i++ {
		h.wg.Done()
	}
}

// AddWait 添加任务协程数
func (h *Handler) AddWait() {
	h.waitCount = h.waitCount + 1
	h.wg.Add(1)
}

// WaitDone 等待任务结束
func (h *Handler) WaitDone() {
	h.wg.Wait()
}

func NewHandler(ctx context.Context) *Handler {
	return &Handler{
		ctx:       ctx,
		waitCount: 0,
		queue:     collections.NewQueue(),
	}
}
