package ack

import (
	"errors"
	"github.com/topcloudz/fvpn/pkg/option"
	"github.com/topcloudz/fvpn/pkg/packet"
	"github.com/topcloudz/fvpn/pkg/packet/header"
	"net"
	"unsafe"
)

type PingPacketAck struct {
	header  header.Header
	IP      net.IP
	Port    uint16
	NatIP   net.IP
	NatPort uint16
}

func NewPacket() PingPacketAck {
	h, _ := header.NewHeader(option.MsgTypePingAck, "")
	return PingPacketAck{header: h}
}

// Encode encode a NotifyPacket to bytes  sequence ip-> port->nattype->destAddr ->natip->natport
func Encode(p PingPacketAck) ([]byte, error) {
	b := make([]byte, unsafe.Sizeof(PingPacketAck{}))
	headerBuff, err := header.Encode(p.header)
	if err != nil {
		return nil, errors.New("encode common packet failed")
	}
	idx := 0
	idx = packet.EncodeBytes(b, headerBuff, idx)
	idx = packet.EncodeBytes(b, p.IP, idx)
	idx = packet.EncodeUint16(b, p.Port, idx)
	idx = packet.EncodeBytes(b, p.NatIP, idx)
	idx = packet.EncodeUint16(b, p.NatPort, idx)

	return b, nil
}

// Decode decode buff to NotifyPacket   ip-> port->nattype->destAddr ->natip->natport
func Decode(buff []byte) (PingPacketAck, error) {
	res := NewPacket()
	h, err := header.Decode(buff)
	if err != nil {
		return PingPacketAck{}, errors.New("decode common packet failed")
	}
	idx := 0
	res.header = h
	idx += int(unsafe.Sizeof(header.Header{}))

	var ip = make([]byte, 16)
	idx = packet.DecodeBytes(&ip, buff, idx)
	res.IP = ip

	idx = packet.DecodeUint16(&res.Port, buff, idx)

	var natIp = make([]byte, 16)
	idx = packet.DecodeBytes(&natIp, buff, idx)
	res.NatIP = natIp
	idx = packet.DecodeUint16(&res.NatPort, buff, idx)

	return res, nil
}
