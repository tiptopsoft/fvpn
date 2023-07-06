package packet

import (
	"fmt"
	"net"
	"testing"
)

func TestEncode(t *testing.T) {
	h, _ := NewHeader(3, "123444444444abcdef")
	h.SrcIP = net.ParseIP("5.244.24.141")
	h.DstIP = net.ParseIP("192.168.0.1")
	buff, _ := Encode(h)
	fmt.Println("len: ", len(buff))
	fmt.Println(buff)

	h1, _ := Decode(buff)
	fmt.Println(h1.SrcIP)
	fmt.Println(h1.DstIP)

}
