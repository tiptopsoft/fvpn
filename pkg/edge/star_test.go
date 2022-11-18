package edge

import (
	"net"
	"testing"
)

func TestUdpRequest(t *testing.T) {

	conn, err := net.Dial("udp", "127.0.0.1:4000")
	if err != nil {
		panic(err)
	}

	//p := common.NewPacket()
	//data, _ := p.Encode()
	if _, err := conn.Write([]byte{1}); err != nil {
		panic(err)
	}

	conn.Close()
}
