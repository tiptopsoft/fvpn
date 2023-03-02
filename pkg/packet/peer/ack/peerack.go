package ack

import (
	"github.com/interstellar-cloud/star/pkg/option"
	packet2 "github.com/interstellar-cloud/star/pkg/packet"
	"github.com/interstellar-cloud/star/pkg/packet/common"
	common2 "github.com/interstellar-cloud/star/pkg/packet/common"
	"net"
	"unsafe"
)

// EdgeInfo info need to connect to
type EdgeInfo struct {
	Mac  net.HardwareAddr
	Host net.IP
	Port uint16
	P2P  uint8 //1: 是2：否
}

// EdgePacketAck ack for size of EdgeInfo
type EdgePacketAck struct {
	common2.CommonPacket
	Size      uint8
	PeerInfos []EdgeInfo
}

func NewPacket() EdgePacketAck {
	cmPacket := common2.NewPacket(option.MsgTypeQueryPeer)
	return EdgePacketAck{
		CommonPacket: cmPacket,
	}
}

func (ack EdgePacketAck) Encode() ([]byte, error) {
	b := make([]byte, 2048)
	cp, err := ack.CommonPacket.Encode()
	if err != nil {
		return nil, err
	}

	idx := 0
	idx = packet2.EncodeBytes(b, cp, idx)
	idx = packet2.EncodeUint8(b, ack.Size, idx)
	for _, v := range ack.PeerInfos {
		idx = packet2.EncodeBytes(b, v.Mac, idx)
		idx = packet2.EncodeBytes(b, v.Host, idx)
		idx = packet2.EncodeUint16(b, v.Port, idx)
	}

	return b, nil
}

func (ack EdgePacketAck) Decode(udpBytes []byte) (packet2.Interface, error) {
	idx := 0
	cp, err := common.NewPacketWithoutType().Decode(udpBytes)
	idx += int(unsafe.Sizeof(common2.CommonPacket{}))
	ack.CommonPacket = cp.(common.CommonPacket)

	idx = packet2.DecodeUint8(&ack.Size, udpBytes, idx)

	var info []EdgeInfo
	for i := 0; uint8(i) < ack.Size; i++ {
		peer := EdgeInfo{}
		var mac = make([]byte, 6)
		idx = packet2.DecodeBytes(&mac, udpBytes, idx)
		peer.Mac = mac
		var ip = make([]byte, 16)
		idx = packet2.DecodeBytes(&ip, udpBytes, idx)
		peer.Host = ip
		idx = packet2.DecodeUint16(&peer.Port, udpBytes, idx)
		info = append(info, peer)
	}

	ack.PeerInfos = info

	if err != nil {
		return ack, err
	}

	return ack, nil
}
