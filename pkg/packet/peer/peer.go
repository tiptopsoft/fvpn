package peer

import (
	"errors"
	"github.com/topcloudz/fvpn/pkg/packet"
	"github.com/topcloudz/fvpn/pkg/util"
	"unsafe"
)

type PeerPacket struct {
	header packet.Header
}

func (pkt PeerPacket) String() string {
	//value := fmt.Sprintf("type: %d, srcMac: %s", pkt.header.Flags, pkt.SrcMac.String())
	//return value
	return ""
}

func NewPacket(networkId string) PeerPacket {
	cmPacket, _ := packet.NewHeader(util.MsgTypeQueryPeer, networkId)
	return PeerPacket{
		header: cmPacket,
	}
}

func Encode(p PeerPacket) ([]byte, error) {
	b := make([]byte, unsafe.Sizeof(PeerPacket{}))
	headerBuff, err := packet.Encode(p.header)
	if err != nil {
		return nil, errors.New("encode common packet failed")
	}
	idx := 0
	idx = packet.EncodeBytes(b, headerBuff, idx)
	//idx = packet.EncodeBytes(b, p.SrcMac[:], idx)
	return b, nil
}

func (p PeerPacket) Decode(udpBytes []byte) (PeerPacket, error) {

	res := NewPacket("")
	h, err := packet.Decode(udpBytes)
	if err != nil {
		return PeerPacket{}, errors.New("decode common packet failed")
	}
	idx := 0
	res.header = h
	idx += int(unsafe.Sizeof(packet.Header{}))
	//var mac = make([]byte, 6)
	//packet.DecodeBytes(&mac, udpBytes, idx)
	//res.SrcMac = mac
	return res, nil
}

func DecodeWithCommonPacket(udpBytes []byte, cp packet.Header) (PeerPacket, error) {
	res := NewPacket("")
	idx := 0
	res.header = cp
	idx += int(unsafe.Sizeof(packet.Header{}))
	//var mac = make([]byte, 6)
	//packet.DecodeBytes(&mac, udpBytes, idx)
	//res.SrcMac = mac
	return res, nil
}
