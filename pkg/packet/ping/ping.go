package ping

import (
	"errors"
	"github.com/topcloudz/fvpn/pkg/option"
	"github.com/topcloudz/fvpn/pkg/packet"
	"github.com/topcloudz/fvpn/pkg/packet/header"
	"net"
	"unsafe"
)

type PingPacket struct {
	header header.Header
	IP     net.IP
	DstIP  net.IP
}

func NewPacket() PingPacket {
	h, _ := header.NewHeader(option.MsgTypePing, "")
	return PingPacket{
		header: h,
	}
}

// Encode encode a NotifyPacket to bytes  sequence ip-> port->nattype->destAddr ->natip->natport
func Encode(p PingPacket) ([]byte, error) {
	b := make([]byte, unsafe.Sizeof(PingPacket{}))
	headerBuff, err := header.Encode(p.header)
	if err != nil {
		return nil, errors.New("encode common packet failed")
	}
	idx := 0
	idx = packet.EncodeBytes(b, headerBuff, idx)
	idx = packet.EncodeBytes(b, p.IP, idx)

	return b, nil
}

// Decode decode buff to NotifyPacket   ip-> port->nattype->destAddr ->natip->natport
func Decode(buff []byte) (PingPacket, error) {
	res := NewPacket()
	h, err := header.Decode(buff)
	if err != nil {
		return PingPacket{}, errors.New("decode common packet failed")
	}
	idx := 0
	res.header = h
	idx += int(unsafe.Sizeof(header.Header{}))

	var ip = make([]byte, 16)
	idx = packet.DecodeBytes(&ip, buff, idx)
	res.IP = ip

	return res, nil
}
