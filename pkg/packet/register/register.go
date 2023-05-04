package register

import (
	"errors"
	"github.com/topcloudz/fvpn/pkg/option"
	"github.com/topcloudz/fvpn/pkg/packet"
	"github.com/topcloudz/fvpn/pkg/packet/header"
	"net"
	"unsafe"
)

// RegPacket server a client to server
type RegPacket struct { //48
	header header.Header    //12
	SrcMac net.HardwareAddr //6
	SrcIP  net.IP           // 4 byte是ipv4, 16 byte是ipv6
}

func NewPacket(networkId string, srcMac net.HardwareAddr, srcIP net.IP) RegPacket {
	cmPacket, _ := header.NewHeader(option.MsgTypeRegisterSuper, networkId)
	reg := RegPacket{
		header: cmPacket,
		SrcIP:  srcIP,
		SrcMac: srcMac,
	}

	return reg
}

func NewUnregisterPacket(networkId string) RegPacket {
	cmPacket, _ := header.NewHeader(option.MsgTypeUnregisterSuper, networkId)
	return RegPacket{
		header: cmPacket,
	}
}

func Encode(regPacket RegPacket) ([]byte, error) {
	b := make([]byte, 48)
	headerBuff, err := header.Encode(regPacket.header)
	if err != nil {
		return nil, errors.New("encode Header failed")
	}
	idx := 0
	idx = packet.EncodeBytes(b, headerBuff, idx)
	idx = packet.EncodeBytes(b, regPacket.SrcMac[:], idx)
	idx = packet.EncodeBytes(b, regPacket.SrcIP[:], idx)
	return b, nil
}

func Decode(buff []byte) (RegPacket, error) {
	res := NewPacket("", net.HardwareAddr{}, net.IP{})

	h, err := header.Decode(buff[:12])
	if err != nil {
		return RegPacket{}, err
	}
	res.header = h
	idx := 0
	idx += int(unsafe.Sizeof(header.Header{}))
	var mac = make([]byte, 6)
	packet.DecodeBytes(&mac, buff, idx)
	copy(res.SrcMac[:], mac)
	var ip = make([]byte, 16)
	packet.DecodeBytes(&ip, buff, idx)
	copy(res.SrcIP[:], ip)
	return res, nil
}
