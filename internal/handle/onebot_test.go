package handle

import (
	"fmt"
	"testing"
)

func TestType(t *testing.T) {
	var delay interface{}
	delay = 200000000000000
	if v, ok := delay.(int); ok {
		fmt.Println("asdsdasdasd", v)
		return
	}
}
