package peer

import (
	"fmt"
	"net"
	"testing"
)

func TestEncode1(t *testing.T) {

	p := NewPeerPacket()
	p.Header.SrcIP = net.ParseIP("121.1.1.1")
	buff, err := Encode(p)
	if err != nil {
		panic(buff)
	}

	p1, _ := Decode(buff)
	fmt.Println(p1.Header.SrcIP)
}
