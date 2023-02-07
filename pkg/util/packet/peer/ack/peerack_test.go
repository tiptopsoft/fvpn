package ack

import (
	"fmt"
	"github.com/interstellar-cloud/star/pkg/util"
	"github.com/magiconair/properties/assert"
	"net"
	"testing"
)

func TestDecode(t *testing.T) {
	mac, _ := util.RandMac()
	fmt.Println([]byte(mac))
	fmt.Println(mac)
	fmt.Println(mac)
	var result []EdgeInfo
	srcMac, err := net.ParseMAC(mac)
	fmt.Println("src: ", srcMac)
	if err != nil {
		panic(err)
	}
	info := EdgeInfo{
		Mac:  srcMac,
		Host: net.IPv4(127, 0, 0, 1),
		Port: 35582,
	}

	result = append(result, info)

	peerPacket := NewPacket()
	peerPacket.Size = 1
	peerPacket.PeerInfos = result

	fmt.Println("origin: ", peerPacket)
	data, _ := Encode(peerPacket)

	fmt.Println("encode code: ", data)

	pack, _ := Decode(data)
	fmt.Println("decode data: ", pack, pack.Size)

	assert.Equal(t, pack.Size, peerPacket.Size)

}
