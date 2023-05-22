package ack

import (
	"fmt"
	"net"
	"testing"
)

func TestEncode(t *testing.T) {

}
func TestDecode(t *testing.T) {

	pkt := NewPacket()
	pkt.Size = 1

	mac, _ := net.ParseMAC("01:e1:22:23:12:12")
	ip := net.ParseIP("127.0.0.1")
	info := EdgeInfo{
		Mac:     mac,
		IP:      ip,
		Port:    4000,
		P2P:     0,
		NatIp:   ip,
		NatPort: 54230,
	}

	pkt.NodeInfos = append(pkt.NodeInfos, info)

	originBuff, _ := Encode(pkt)

	fmt.Println("ori packet:", originBuff)

	ac, err := Decode(originBuff)
	fmt.Println(ac, err)

}
