package peer

import (
	"errors"
	"github.com/interstellar-cloud/star/pkg/option"
	"github.com/interstellar-cloud/star/pkg/packet"
	"github.com/interstellar-cloud/star/pkg/packet/common"
	"net"
	"unsafe"
)

type PeerPacket struct {
	header common.PacketHeader
	SrcMac net.HardwareAddr
}

func NewPacket() PeerPacket {
	cmPacket := common.NewPacket(option.MsgTypeQueryPeer)
	return PeerPacket{
		header: cmPacket,
	}
}

func (p PeerPacket) Encode() ([]byte, error) {
	b := make([]byte, unsafe.Sizeof(PeerPacket{}))
	commonBytes, err := p.header.Encode()
	if err != nil {
		return nil, errors.New("encode common packet failed")
	}
	idx := 0
	idx = packet.EncodeBytes(b, commonBytes, idx)
	idx = packet.EncodeBytes(b, p.SrcMac[:], idx)
	return b, nil
}

func (p PeerPacket) Decode(udpBytes []byte) (packet.Interface, error) {

	res := NewPacket()
	cp, err := common.NewPacketWithoutType().Decode(udpBytes)
	if err != nil {
		return PeerPacket{}, errors.New("decode common packet failed")
	}
	idx := 0
	res.header = cp.(common.PacketHeader)
	idx += int(unsafe.Sizeof(common.PacketHeader{}))
	var mac = make([]byte, 6)
	packet.DecodeBytes(&mac, udpBytes, idx)
	res.SrcMac = mac
	return res, nil
}

func DecodeWithCommonPacket(udpBytes []byte, cp common.PacketHeader) (PeerPacket, error) {
	res := NewPacket()
	idx := 0
	res.header = cp
	idx += int(unsafe.Sizeof(common.PacketHeader{}))
	var mac = make([]byte, 6)
	packet.DecodeBytes(&mac, udpBytes, idx)
	res.SrcMac = mac
	return res, nil
}
