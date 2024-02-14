package handle

import (
	"encoding/json"
	"fmt"
	global "github.com/IUnlimit/perpetua/internal"
	"github.com/IUnlimit/perpetua/internal/conf"
	"github.com/IUnlimit/perpetua/internal/utils"
	"testing"
)

func TestCreateWSInstance(t *testing.T) {
	conf.Init()
	//CreateWSInstance(25565)
}

func TestUnmarshall(t *testing.T) {
	msg := "{\"interval\":5000,\"status\":{\"app_initialized\":true,\"app_enabled\":true,\"app_good\":true,\"online\":true,\"good\":true},\"meta_event_type\":\"heartbeat\",\"time\":1705851315,\"self_id\":3012218237,\"post_type\":\"meta_event\"}"
	var event global.MsgData
	_ = json.Unmarshal([]byte(msg), &event)
	//fmt.Println(event)

	_ = event["status"].(global.MsgData)
	//fmt.Println(status)
}

func TestBuild(t *testing.T) {
	response := utils.BuildWSResponse("ok", 0, "echo",
		"uuid", "7657")
	fmt.Println(response)
}
