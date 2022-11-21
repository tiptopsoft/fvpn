package ack

import (
	"github.com/interstellar-cloud/star/pkg/packet"
	"github.com/interstellar-cloud/star/pkg/packet/common"
	"net"
	"unsafe"
)

// PeerInfo info need to connect to
type PeerInfo struct {
	Mac  net.HardwareAddr
	Host net.IP
	Port uint16
	P2p  uint8 //1: 是2：否
}

// PeerPacketAck ack for size of PeerInfo
type PeerPacketAck struct {
	common.CommonPacket
	Size      uint8
	PeerInfos []PeerInfo
}

func NewPacket() PeerPacketAck {
	return PeerPacketAck{}
}

func Encode(ack PeerPacketAck) ([]byte, error) {
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

func Decode(udpBytes []byte) (PeerPacketAck, error) {
	ack := PeerPacketAck{}
	idx := 0
	cp, err := common.Decode(udpBytes)
	idx += int(unsafe.Sizeof(common.NewPacket()))
	ack.CommonPacket = cp

	idx = packet.DecodeUint8(&ack.Size, udpBytes, idx)
	idx = packet.DecodeUint8(&ack.Size, udpBytes, idx)

	var info []PeerInfo
	for i := 0; uint8(i) < ack.Size; i++ {
		peer := PeerInfo{}
		var mac = make([]byte, 6)
		idx = packet.DecodeBytes(&mac, udpBytes, idx)
		peer.Mac = mac
		var ip = make([]byte, 4)
		idx = packet.DecodeBytes(&ip, udpBytes, idx)
		idx = packet.DecodeUint16(&peer.Port, udpBytes, idx)
		info = append(info, peer)
	}

	ack.PeerInfos = info

	if err != nil {
		return ack, err
	}

	return ack, nil
}
