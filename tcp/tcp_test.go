package tcp

import (
	"fmt"
	"testing"
)

var (
	data1 OffsetControlFlag = 0xa002
)

func TestOffsetControlFlag(t *testing.T) {
	fmt.Printf("offset: %d\n", data1.Offset())
	fmt.Println("flag: ", data1.ControlFlag().String())
}
