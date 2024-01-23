package handle

import (
	"context"
	"sync"
)

type Handler struct {
	ctx       context.Context
	wg        sync.WaitGroup
	cache     [][]byte
	waitCount int
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
	}
}
