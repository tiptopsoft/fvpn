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
	info := EdgeInfo{
		Mac:     mac,
		IP:      net.IPv4(127, 0, 0, 1),
		Port:    4000,
		P2P:     0,
		NatIp:   net.IP{127, 0, 0, 1},
		NatPort: 54230,
	}

	pkt.NodeInfos = append(pkt.NodeInfos, info)

	originBuff, _ := Encode(pkt)

	fmt.Println("ori packet:", originBuff)

	ac, err := Decode(originBuff)
	fmt.Println(ac, err)
	//mac, _ := util.RandMac()
	//fmt.Println([]byte(mac))
	//fmt.Println(mac)
	//fmt.Println(mac)
	//var result []EdgeInfo
	//srcMac, err := net.ParseMAC(mac)
	//fmt.Println("src: ", srcMac)
	//if err != nil {
	//	panic(err)
	//}
	//info := EdgeInfo{
	//	Mac:  srcMac,
	//	IP: net.IPv4(127, 0, 0, 1),
	//	Port: 35582,
	//}
	//
	//result = append(result, info)
	//
	//peerPacket := NewPacket()
	//peerPacket.Size = 1
	//peerPacket.NodeInfos = result
	//
	//fmt.Println("origin: ", peerPacket)
	//data, _ := Encode(peerPacket)
	//
	//fmt.Println("encode code: ", data)
	//
	//pack, _ := Decode(data)
	//fmt.Println("decode data: ", pack, pack.Size)
	//
	//assert.Equal(t, pack.Size, peerPacket.Size)

}
