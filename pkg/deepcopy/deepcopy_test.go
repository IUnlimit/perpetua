package deepcopy

import (
	"encoding/json"
	"fmt"
	global "github.com/IUnlimit/perpetua/internal"
	"testing"
)

func TestCopy(t *testing.T) {
	msg := "{\"interval\":5000,\"status\":{\"app_initialized\":true,\"app_enabled\":true,\"app_good\":true,\"online\":true,\"good\":true},\"meta_event_type\":\"heartbeat\",\"time\":1705851315,\"self_id\":3012218237,\"post_type\":\"meta_event\"}"
	var event *global.MsgData
	_ = json.Unmarshal([]byte(msg), &event)

	newMsg := Copy(event)
	newStatus := (*(newMsg.(*global.MsgData)))["status"]
	newStatus.(global.MsgData)["app_initialized"] = false

	fmt.Println("Ori: ", (*event)["status"])
	fmt.Println("Now: ", newStatus)
}
