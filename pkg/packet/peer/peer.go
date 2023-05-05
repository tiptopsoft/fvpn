package peer

import (
	"errors"
	"github.com/topcloudz/fvpn/pkg/option"
	"github.com/topcloudz/fvpn/pkg/packet"
	"github.com/topcloudz/fvpn/pkg/packet/header"
	"net"
	"unsafe"
)

type PeerPacket struct {
	header header.Header
	SrcMac net.HardwareAddr
}

func NewPacket(networkId string) PeerPacket {
	cmPacket, _ := header.NewHeader(option.MsgTypeQueryPeer, networkId)
	return PeerPacket{
		header: cmPacket,
	}
}

func Encode(p PeerPacket) ([]byte, error) {
	b := make([]byte, unsafe.Sizeof(PeerPacket{}))
	headerBuff, err := header.Encode(p.header)
	if err != nil {
		return nil, errors.New("encode common packet failed")
	}
	idx := 0
	idx = packet.EncodeBytes(b, headerBuff, idx)
	idx = packet.EncodeBytes(b, p.SrcMac[:], idx)
	return b, nil
}

func (p PeerPacket) Decode(udpBytes []byte) (PeerPacket, error) {

	res := NewPacket("")
	h, err := header.Decode(udpBytes)
	if err != nil {
		return PeerPacket{}, errors.New("decode common packet failed")
	}
	idx := 0
	res.header = h
	idx += int(unsafe.Sizeof(header.Header{}))
	var mac = make([]byte, 6)
	packet.DecodeBytes(&mac, udpBytes, idx)
	res.SrcMac = mac
	return res, nil
}

func DecodeWithCommonPacket(udpBytes []byte, cp header.Header) (PeerPacket, error) {
	res := NewPacket("")
	idx := 0
	res.header = cp
	idx += int(unsafe.Sizeof(header.Header{}))
	var mac = make([]byte, 6)
	packet.DecodeBytes(&mac, udpBytes, idx)
	res.SrcMac = mac
	return res, nil
}
