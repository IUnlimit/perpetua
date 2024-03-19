package handle

import (
	global "github.com/IUnlimit/perpetua/internal"
	collections "github.com/chenjiandongx/go-queue"
)

// todo replace by block queue
var echoMap *EchoMap

// EchoMap 利用 echo 完成 NTQQ -> client 消息调度 (发送时标记echo)
type EchoMap struct {
	Receive chan bool

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
	queue := echoMap.dataMap[id]
	if queue == nil {
		queue = collections.NewQueue()
		echoMap.dataMap[id] = queue
	}
	queue.Put(data)
}

func (em *EchoMap) JustGet(id string, consumer func(global.MsgData)) {
	queue := echoMap.dataMap[id]
	if queue == nil {
		return
	}

	for {
		e, ok := queue.Get()
		if !ok {
			return
		}
		consumer(e.(global.MsgData))
	}
}
