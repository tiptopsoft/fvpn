package ack

import (
	"fmt"
	"testing"
	"unsafe"
)

func TestNewPacket(t *testing.T) {
	fmt.Println(unsafe.Sizeof(RegPacketAck{}))
}
