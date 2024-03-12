package handle

import (
	"fmt"
	"testing"
)

func TestTryReverseWSInstance(t *testing.T) {
	err := TryReverseWSInstance("ws://127.0.0.1:8800/onebot/v11/ws", "")
	if err != nil {
		fmt.Printf("%v", err)
	}
}
