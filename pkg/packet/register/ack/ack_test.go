package ack

import (
	"fmt"
	"github.com/topcloudz/fvpn/pkg/packet/header"
	"net"
	"testing"
	"unsafe"

	"github.com/magiconair/properties/assert"
	"github.com/topcloudz/fvpn/pkg/option"
)

func TestNewPacket(t *testing.T) {

	size := unsafe.Sizeof(header.Header{})
	fmt.Println(size)

	RecMac := "01:01:03:02:03:01"
	ip := "192.168.1.1"
	Mask := "255.255.255.0"

	p := NewPacket()
	mac, err := net.ParseMAC(RecMac)
	if err != nil {
		panic(err)
	}

	p.RegMac = mac
	vip := net.ParseIP(ip)
	ipsize := unsafe.Sizeof(vip)
	fmt.Println(ipsize)
	p.AutoIP = vip
	p.Mask = net.ParseIP(Mask)
	p.header = header.NewHeader(option.MsgTypeRegisterAck)

	buff, err := p.Encode()
	fmt.Println(buff)

	//decod
	p1 := NewPacket()
	res, err := p1.Decode(buff)
	if err != nil {
		panic(err)
	}

	r := res.(RegPacketAck)
	assert.Equal(t, r.AutoIP.String(), ip)
}
