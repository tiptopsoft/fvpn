package ack

import (
	"github.com/topcloudz/fvpn/pkg/packet"
	"github.com/topcloudz/fvpn/pkg/util"
	"net"
	"unsafe"
)

// EdgeInfo info need to connect to
type EdgeInfo struct {
	Mac     net.HardwareAddr
	IP      net.IP
	Port    uint16
	P2P     uint8 //1: 是2：否
	NatIp   net.IP
	NatPort uint16
}

// EdgePacketAck ack for size of EdgeInfo
type EdgePacketAck struct {
	header    packet.Header
	Size      uint8
	NodeInfos []EdgeInfo
}

func NewPacket(networkId string) EdgePacketAck {
	cmPacket, _ := packet.NewHeader(util.MsgTypeQueryPeer, networkId)
	return EdgePacketAck{
		header: cmPacket,
	}
}

func Encode(ack EdgePacketAck) ([]byte, error) {
	length := 12 + 1 + ack.Size*42
	b := make([]byte, length)
	cp, err := packet.Encode(ack.header)
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
		idx = packet.EncodeBytes(b, v.NatIp, idx)
		idx = packet.EncodeUint16(b, v.NatPort, idx)
	}

	return b, nil
}

func Decode(udpBytes []byte) (EdgePacketAck, error) {
	ack := NewPacket("")
	idx := 0
	h, err := packet.Decode(udpBytes)
	idx += int(unsafe.Sizeof(packet.Header{}))
	ack.header = h

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
		var natIp = make([]byte, 16)
		idx = packet.DecodeBytes(&natIp, udpBytes, idx)
		peer.NatIp = natIp
		idx = packet.DecodeUint16(&peer.NatPort, udpBytes, idx)
		info = append(info, peer)
	}

	ack.NodeInfos = info

	if err != nil {
		return ack, err
	}

	return ack, nil
}
