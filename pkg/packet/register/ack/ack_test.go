package ack

import (
	"fmt"
	"github.com/interstellar-cloud/star/pkg/option"
	"github.com/interstellar-cloud/star/pkg/packet/common"
	"net"
	"testing"
	"unsafe"
)

func TestNewPacket(t *testing.T) {

	size := unsafe.Sizeof(common.CommonPacket{})
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
	p.CommonPacket = common.NewPacket(option.MsgTypeRegisterAck)

	fmt.Println(Encode(p))
}
