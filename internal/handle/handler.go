package handle

import (
	"context"
	global "github.com/IUnlimit/perpetua/internal"
	"github.com/IUnlimit/perpetua/pkg/concurrent"
	collections "github.com/chenjiandongx/go-queue"
	"github.com/google/uuid"
	"sync"
)

// handleSet stores handle.Handler for client websocket
var handleSet = concurrent.NewSet()

type Handler struct {
	Receive chan bool
	Lock    sync.Mutex

	id string
	// TODO
	// nullable, can only be actively set
	name string
	ctx  context.Context
	// waiting goroutine count
	waitCount int
	// thread safe queue
	queue *collections.Queue
	wg    sync.WaitGroup
}

func (h *Handler) AddMessage(uuid string) {
	h.queue.Put(uuid)
}

// GetMessage from local cache
func (h *Handler) GetMessage(consumer func(data global.MsgData)) {
	for {
		e, ok := h.queue.Get()
		if !ok {
			return
		}
		data, _ := globalCache.cache.Get(e)
		if data == nil {
			continue
		}
		consumer(data.(global.MsgData))
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
	h.waitCount = 0
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

func (h *Handler) GetId() string {
	return h.id
}

func (h *Handler) GetName() string {
	return h.name
}

func NewHandler(ctx context.Context) *Handler {
	return &Handler{
		ctx:       ctx,
		id:        uuid.NewString(),
		Receive:   make(chan bool),
		waitCount: 0,
		queue:     collections.NewQueue(),
	}
}

func FindHandler(id string) *Handler {
	for _, v := range handleSet.Iterator() {
		handler := v.(*Handler)
		if handler.id == id {
			return handler
		}
	}
	return nil
}
