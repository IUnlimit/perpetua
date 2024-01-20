package utils

import (
	"runtime"
	"strings"
)

func IsWinPlatform() bool {
	return strings.Contains(strings.ToLower(runtime.GOOS), "win")
}
