package handle

import (
	"sync"

	global "github.com/IUnlimit/perpetua/internal"
	collections "github.com/chenjiandongx/go-queue"
)

// todo replace by block queue
var echoMap *EchoMap

// EchoMap 利用 echo 完成 NTQQ -> client 消息调度 (发送时标记echo)
type EchoMap struct {
	Receive chan bool

	mu      sync.Mutex
	dataMap map[string]*collections.Queue
}

func NewEchoMap() *EchoMap {
	return &EchoMap{
		dataMap: make(map[string]*collections.Queue),
		Receive: make(chan bool),
	}
}

// JustPut echo
func (em *EchoMap) JustPut(id string, data global.MsgData) {
	em.mu.Lock()
	queue := em.dataMap[id]
	if queue == nil {
		queue = collections.NewQueue()
		em.dataMap[id] = queue
	}
	em.mu.Unlock()
	queue.Put(data)
}

// JustGet drains queued data of id.
// canContinue is checked before taking each element, returning false stops draining and keeps the remaining data queued.
func (em *EchoMap) JustGet(id string, canContinue func() bool, consumer func(global.MsgData)) {
	em.mu.Lock()
	queue := em.dataMap[id]
	em.mu.Unlock()
	if queue == nil {
		return
	}

	for {
		if canContinue != nil && !canContinue() {
			return
		}
		e, ok := queue.Get()
		if !ok {
			return
		}
		consumer(e.(global.MsgData))
	}
}
