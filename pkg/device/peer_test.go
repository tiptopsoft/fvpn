package device

import (
	"fmt"
	"net"
	"testing"
)

func TestPeer_GetIP(t *testing.T) {
	s := "tiptopsoft.cn"
	addr, err := net.ResolveUDPAddr("udp", s)
	if err != nil {
		t.Fatal(err)
	}

	fmt.Println(addr.IP)
}
