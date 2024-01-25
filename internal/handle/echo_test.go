package handle

import (
	"context"
	"encoding/json"
	"fmt"
	global "github.com/IUnlimit/perpetua/internal"
	collections "github.com/chenjiandongx/go-queue"
	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
	"testing"
)

func TestEcho(t *testing.T) {
	id := uuid.NewString()
	str := fmt.Sprintf("%s#%s#a", global.EchoPrefix, id)
	re := global.EchoRegx
	matches := re.FindStringSubmatch(str)

	if len(matches) == 4 {
		group1 := matches[1]
		group2 := matches[2]
		group3 := matches[3]

		fmt.Println("Group 1:", group1)
		fmt.Println("Group 2:", group2)
		fmt.Println("Group 3:", group3)
	} else {
		fmt.Println("No match found, len: ", len(uuid.NewString()))
	}
}

func TestUpdateEcho(t *testing.T) {
	message := `{
    "action": "send_private_msg",
    "params": {
        "user_id": 765743073,
        "message": "hello"
    },
    "echo": ""
}`
	var msgData global.MsgData
	err := json.Unmarshal([]byte(message), &msgData)
	if err != nil {
		log.Errorf("[<-Client] Failed to unmarshal client message: %s", string(message))
		return
	}

	// sign with echo field
	id := uuid.NewString()
	echo := msgData["echo"].(string)
	echo = fmt.Sprintf("%s#%s#%s", global.EchoPrefix, id, echo)
	msgData["echo"] = echo
	log.Printf("[<-Client] Update client(port-%d) message echo: %s", 666, echo)
}

func TestSend(t *testing.T) {
	message := `{
    "action": "send_private_msg",
    "params": {
        "user_id": 765743073,
        "message": "hello"
    },
    "echo": ""
}`
	var msgData global.MsgData
	err := json.Unmarshal([]byte(message), &msgData)
	if err != nil {
		log.Errorf("[<-Client] Failed to unmarshal client message: %s", string(message))
		return
	}

	bytes, _ := json.Marshal(msgData)
	fmt.Println(string(bytes))
}

func TestQueue(t *testing.T) {
	queue := collections.NewQueue()
	queue.Put("666")
	for {
		e, ok := queue.Get()
		if !ok {
			return
		}
		fmt.Println(e)
	}
}

func TestSlice(t *testing.T) {
	receivers := make([]*Handler, 1)
	handler := NewHandler(context.Background())
	receivers = append(receivers, handler)
	fmt.Println(len(receivers))
}
