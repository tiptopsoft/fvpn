package ack

import (
	"github.com/interstellar-cloud/star/pkg/util/option"
	"github.com/interstellar-cloud/star/pkg/util/packet"
	"github.com/interstellar-cloud/star/pkg/util/packet/common"
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
	common.CommonPacket
	Size      uint8
	PeerInfos []EdgeInfo
}

func NewPacket() EdgePacketAck {
	cmPacket := common.NewPacket(option.MsgTypeQueryPeer)
	return EdgePacketAck{
		CommonPacket: cmPacket,
	}
}

func Encode(ack EdgePacketAck) ([]byte, error) {
	b := make([]byte, 2048)
	cp, err := common.Encode(ack.CommonPacket)
	if err != nil {
		return nil, err
	}

	idx := 0
	idx = packet.EncodeBytes(b, cp, idx)
	idx = packet.EncodeUint8(b, ack.Size, idx)
	for _, v := range ack.PeerInfos {
		idx = packet.EncodeBytes(b, v.Mac, idx)
		idx = packet.EncodeBytes(b, v.Host, idx)
		idx = packet.EncodeUint16(b, v.Port, idx)
	}

	return b, nil
}

func Decode(udpBytes []byte) (EdgePacketAck, error) {
	ack := EdgePacketAck{}
	idx := 0
	cp, err := common.Decode(udpBytes)
	idx += int(unsafe.Sizeof(common.CommonPacket{}))
	ack.CommonPacket = cp

	idx = packet.DecodeUint8(&ack.Size, udpBytes, idx)

	var info []EdgeInfo
	for i := 0; uint8(i) < ack.Size; i++ {
		peer := EdgeInfo{}
		var mac = make([]byte, 6)
		idx = packet.DecodeBytes(&mac, udpBytes, idx)
		peer.Mac = mac
		var ip = make([]byte, 16)
		idx = packet.DecodeBytes(&ip, udpBytes, idx)
		peer.Host = ip
		idx = packet.DecodeUint16(&peer.Port, udpBytes, idx)
		info = append(info, peer)
	}

	ack.PeerInfos = info

	if err != nil {
		return ack, err
	}

	return ack, nil
}
