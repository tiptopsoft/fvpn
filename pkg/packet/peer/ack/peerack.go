package ack

import (
	"github.com/topcloudz/fvpn/pkg/option"
	"github.com/topcloudz/fvpn/pkg/packet"
	"net"
	"unsafe"
)

// EdgeInfo info need to connect to
type EdgeInfo struct {
	Mac  net.HardwareAddr
	IP   net.IP
	Port uint16
	P2P  uint8 //1: 是2：否
}

// EdgePacketAck ack for size of EdgeInfo
type EdgePacketAck struct {
	header    packet.Header
	Size      uint8
	NodeInfos []EdgeInfo
}

func NewPacket() EdgePacketAck {
	cmPacket := packet.NewHeader(option.MsgTypeQueryPeer, "")
	return EdgePacketAck{
		header: cmPacket,
	}
}

func (ack EdgePacketAck) Encode() ([]byte, error) {
	b := make([]byte, 2048)
	cp, err := ack.header.Encode()
	if err != nil {
		return nil, err
	}

	idx := 0
	idx = packet.EncodeBytes(b, cp, idx)
	idx = packet.EncodeUint8(b, ack.Size, idx)
	for _, v := range ack.NodeInfos {
		idx = packet.EncodeBytes(b, v.Mac, idx)
		idx = packet.EncodeBytes(b, v.IP, idx)
		idx = packet.EncodeUint16(b, v.Port, idx)
	}

	return b, nil
}

func (ack EdgePacketAck) Decode(udpBytes []byte) (packet.Interface, error) {
	idx := 0
	cp, err := packet.NewPacketWithoutType().Decode(udpBytes)
	idx += int(unsafe.Sizeof(packet.Header{}))
	ack.header = cp.(packet.Header)

	idx = packet.DecodeUint8(&ack.Size, udpBytes, idx)

	var info []EdgeInfo
	for i := 0; uint8(i) < ack.Size; i++ {
		peer := EdgeInfo{}
		var mac = make([]byte, 6)
		idx = packet.DecodeBytes(&mac, udpBytes, idx)
		peer.Mac = mac
		var ip = make([]byte, 16)
		idx = packet.DecodeBytes(&ip, udpBytes, idx)
		peer.IP = ip
		idx = packet.DecodeUint16(&peer.Port, udpBytes, idx)
		info = append(info, peer)
	}

	ack.NodeInfos = info

	if err != nil {
		return ack, err
	}

	return ack, nil
}
