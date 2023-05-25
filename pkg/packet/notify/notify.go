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
	header   header.Header //12
	SourceIP net.IP        // inner ip 16
	Port     uint16        // inner port 2
	NatIP    net.IP        // nat ip 16
	NatPort  uint16        //nat port2
	NatType  uint8         //1 retrict 2 symmtrict nat 1
	DestAddr net.IP        //目标IP 2
}

func NewPacket(networkId string) NotifyPacket {
	headerPacket, _ := header.NewHeader(option.MsgTypeNotify, networkId)
	return NotifyPacket{
		header: headerPacket,
	}
}

// Encode encode a NotifyPacket to bytes  sequence ip-> port->nattype->destAddr ->natip->natport
func Encode(np NotifyPacket) ([]byte, error) {
	b := make([]byte, unsafe.Sizeof(NotifyPacket{}))
	headerBuff, err := header.Encode(np.header)
	if err != nil {
		return nil, errors.New("encode common packet failed")
	}
	idx := 0
	idx = packet.EncodeBytes(b, headerBuff, idx)
	idx = packet.EncodeBytes(b, np.SourceIP, idx)
	idx = packet.EncodeUint16(b, np.Port, idx)
	idx = packet.EncodeUint8(b, np.NatType, idx)
	idx = packet.EncodeBytes(b, np.DestAddr, idx)
	idx = packet.EncodeBytes(b, np.NatIP, idx)
	idx = packet.EncodeUint16(b, np.NatPort, idx)

	return b, nil
}

// Decode decode buff to NotifyPacket   ip-> port->nattype->destAddr ->natip->natport
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
	idx = packet.DecodeBytes(&ip, buff, idx)
	res.SourceIP = ip

	idx = packet.DecodeUint16(&res.Port, buff, idx)
	idx = packet.DecodeUint8(&res.NatType, buff, idx)

	var destIp = make([]byte, 16)
	idx = packet.DecodeBytes(&destIp, buff, idx)
	res.DestAddr = destIp

	var natIp = make([]byte, 16)
	idx = packet.DecodeBytes(&natIp, buff, idx)
	res.NatIP = natIp
	idx = packet.DecodeUint16(&res.NatPort, buff, idx)

	return res, nil
}
