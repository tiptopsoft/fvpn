package peer

import (
	"errors"
	"github.com/interstellar-cloud/star/pkg/option"
	packet2 "github.com/interstellar-cloud/star/pkg/packet"
	"github.com/interstellar-cloud/star/pkg/packet/common"
	common2 "github.com/interstellar-cloud/star/pkg/packet/common"
	"net"
	"unsafe"
)

type PeerPacket struct {
	common2.CommonPacket
	SrcMac net.HardwareAddr
}

func NewPacket() PeerPacket {
	cmPacket := common2.NewPacket(option.MsgTypeQueryPeer)
	return PeerPacket{
		CommonPacket: cmPacket,
	}
}

func (p PeerPacket) Encode() ([]byte, error) {
	b := make([]byte, unsafe.Sizeof(PeerPacket{}))
	commonBytes, err := p.CommonPacket.Encode()
	if err != nil {
		return nil, errors.New("encode common packet failed")
	}
	idx := 0
	idx = packet2.EncodeBytes(b, commonBytes, idx)
	idx = packet2.EncodeBytes(b, p.SrcMac[:], idx)
	return b, nil
}

func (p PeerPacket) Decode(udpBytes []byte) (packet2.Interface, error) {

	res := NewPacket()
	cp, err := common.NewPacketWithoutType().Decode(udpBytes)
	if err != nil {
		return PeerPacket{}, errors.New("decode common packet failed")
	}
	idx := 0
	res.CommonPacket = cp.(common.CommonPacket)
	idx += int(unsafe.Sizeof(common2.CommonPacket{}))
	var mac = make([]byte, 6)
	packet2.DecodeBytes(&mac, udpBytes, idx)
	res.SrcMac = mac
	return res, nil
}

func DecodeWithCommonPacket(udpBytes []byte, cp common2.CommonPacket) (PeerPacket, error) {
	res := NewPacket()
	idx := 0
	res.CommonPacket = cp
	idx += int(unsafe.Sizeof(common2.CommonPacket{}))
	var mac = make([]byte, 6)
	packet2.DecodeBytes(&mac, udpBytes, idx)
	res.SrcMac = mac
	return res, nil
}
