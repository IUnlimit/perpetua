package utils

import (
	"errors"
	"fmt"
	global "github.com/IUnlimit/perpetua/internal"
	"github.com/IUnlimit/perpetua/internal/model"
	"github.com/bytedance/gopkg/util/gopool"
	"time"
)

func BuildWSGoodResponse(status string, echo string, entry ...any) global.MsgData {
	if len(entry)%2 != 0 {
		return BuildWSBadResponse(fmt.Sprintf("错误的 map 参数个数: %d", len(entry)), echo)
	}
	if entry == nil {
		return BuildWSResponse(status, 0, echo)
	}
	return BuildWSResponse(status, 0, echo, entry...)
}

func BuildWSBadResponse(status string, echo string) global.MsgData {
	return BuildWSResponse(status, -1, echo)
}

func BuildWSResponse(status string, retcode int32, echo string, entry ...any) global.MsgData {
	m := make(map[string]any)
	for i := 0; i < len(entry); i += 2 {
		m[entry[i].(string)] = entry[i+1]
	}

	return global.MsgData{
		"status":  status,
		"retcode": retcode,
		"data":    m,
		"echo":    echo,
	}
}

// GetForwardImpl get first forwardImpl from appsettings.json
func GetForwardImpl() (*model.Implementation, error) {
	for _, impl := range global.AppSettings.Implementations {
		if impl.Type == "ForwardWebSocket" {
			return impl, nil
		}
	}
	return nil, errors.New("can't find forward websocket impl")
}

// WaitNTQQStartup wait the NTQQ startup
func WaitNTQQStartup(host string, port int, waitCallback func(error)) <-chan struct{} {
	return WaitCondition(time.Duration(1000), func() error {
		return CheckPort(host, port, time.Second*1)
	}, waitCallback)
}

func WaitCondition(gap time.Duration, condition func() error, waitCallback func(error)) <-chan struct{} {
	done := make(chan struct{})

	gopool.Go(func() {
		for {
			err := condition()
			if err == nil {
				break
			}
			if waitCallback != nil {
				waitCallback(err)
			}
			time.Sleep(time.Millisecond * gap)
		}
		close(done)
	})

	return done
}
