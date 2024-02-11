package utils

import (
	"fmt"
	global "github.com/IUnlimit/perpetua/internal"
)

func BuildWSGoodResponse(status string, echo string, entry ...any) global.MsgData {
	if len(entry)%2 != 0 {
		return BuildWSBadResponse(fmt.Sprintf("错误的 map 参数个数: %d", len(entry)), echo)
	}
	if entry == nil {
		return BuildWSResponse(status, 200, echo)
	}
	return BuildWSResponse(status, 200, echo, entry)
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