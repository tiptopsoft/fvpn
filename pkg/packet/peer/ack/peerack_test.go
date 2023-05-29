package ack

import (
	"fmt"
	"net"
	"strconv"
	"strings"
	"testing"
)

func TestEncode(t *testing.T) {

}
func TestDecode(t *testing.T) {

	pkt := NewPacket("")
	pkt.Size = 1

	mac, _ := net.ParseMAC("01:e1:22:23:12:12")
	ip := net.ParseIP("127.0.0.1")
	info := EdgeInfo{
		Mac:     mac,
		IP:      ip,
		Port:    4000,
		P2P:     1,
		NatIp:   ip,
		NatPort: 54230,
	}

	pkt.NodeInfos = append(pkt.NodeInfos, info)

	originBuff, _ := Encode(pkt)

	fmt.Println("ori packet:", originBuff)

	ac, err := Decode(originBuff)
	fmt.Println(ac, err)

}

func TestDecode2(t *testing.T) {
	s := "1 100 0 11 102 46 44 58 123 23 158 255 3 250 52 47 4 154 39 0 0 0 0 0 0 0 0 0 0 255 255 192 168 0 9 0 0 0 0 0 0 0 0 0 0 0 0 255 255 81 70 36 156 232 95 198 96 124 146 141 238 0 0 0 0 0 0 0 0 0 0 255 255 192 168 0 6 0 0 0 0 0 0 0 0 0 0 0 0 255 255 223 108 79 97 165 157 138 125 230 172 19 88 0 0 0 0 0 0 0 0 0 0 255 255 192 168 0 4 0 0 0 0 0 0 0 0 0 0 0 0 255 255 101 43 97 112 23 173"
	arr := strings.Split(s, " ")
	var buff []byte
	for _, v := range arr {
		value, _ := strconv.Atoi(v)
		buff = append(buff, byte(value))
	}

	pack, err := Decode(buff)

	fmt.Println(pack, err)

}
