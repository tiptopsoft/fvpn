package notify

import (
	"errors"
	"github.com/topcloudz/fvpn/pkg/option"
	"github.com/topcloudz/fvpn/pkg/packet"
	"github.com/topcloudz/fvpn/pkg/packet/header"
	"net"
	"unsafe"
)

// NotifyPacket use to tell dest node to connect, punch hole
type NotifyPacket struct {
	header  header.Header
	Addr    net.IP // inner ip
	Port    uint16 // inner port
	NatAddr net.IP // nat ip
	NatPort uint16 //nat port
	NatType uint8  //1 retrict 2 symmtrict nat
}

func NewPacket(networkId string) NotifyPacket {
	headerPacket, _ := header.NewHeader(option.MsgTypeNotify, networkId)
	return NotifyPacket{
		header: headerPacket,
	}
}

// Encode encode a NotifyPacket to bytes
func Encode(np NotifyPacket) ([]byte, error) {
	b := make([]byte, unsafe.Sizeof(NotifyPacket{}))
	headerBuff, err := header.Encode(np.header)
	if err != nil {
		return nil, errors.New("encode common packet failed")
	}
	idx := 0
	idx = packet.EncodeBytes(b, headerBuff, idx)
	idx = packet.EncodeBytes(b, np.Addr, idx)
	idx = packet.EncodeUint16(b, np.Port, idx)
	idx = packet.EncodeBytes(b, np.NatAddr, idx)
	idx = packet.EncodeUint16(b, np.NatPort, idx)
	idx = packet.EncodeUint8(b, np.NatType, idx)
	return b, nil
}

// Decode decode buff to NotifyPacket
func Decode(buff []byte) (NotifyPacket, error) {
	res := NewPacket("")
	h, err := header.Decode(buff)
	if err != nil {
		return NotifyPacket{}, errors.New("decode common packet failed")
	}
	idx := 0
	res.header = h
	idx += int(unsafe.Sizeof(header.Header{}))
	var ip = make([]byte, 16)
	packet.DecodeBytes(&ip, buff, idx)
	res.Addr = ip
	idx = packet.DecodeUint16(&res.Port, buff, idx)
	var natIp = make([]byte, 16)
	idx = packet.DecodeBytes(&natIp, buff, idx)
	res.NatAddr = natIp
	idx = packet.DecodeUint16(&res.NatPort, buff, idx)
	idx = packet.DecodeUint8(&res.NatType, buff, idx)
	return res, nil
}