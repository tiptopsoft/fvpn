package peer

import (
	"fmt"
	"net"
	"testing"
	"unsafe"
)

func TestEncode1(t *testing.T) {
	fmt.Println(unsafe.Sizeof(net.UDPAddr{}))
	fmt.Println(unsafe.Sizeof(net.IP{}))
}
