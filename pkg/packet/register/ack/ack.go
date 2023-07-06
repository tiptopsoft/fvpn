package ack

import (
	packet "github.com/topcloudz/fvpn/pkg/packet"
	"github.com/topcloudz/fvpn/pkg/util"
	"net"
	"unsafe"
)

type RegPacketAck struct {
	header packet.Header    //8 byte
	RegMac net.HardwareAddr //6 byte
	AutoIP net.IP           //4byte
	Mask   net.IP
}

func NewPacket() RegPacketAck {
	cmPacket, _ := packet.NewHeader(util.MsgTypeRegisterAck, "")
	return RegPacketAck{
		header: cmPacket,
	}
}

func Encode(ack RegPacketAck) ([]byte, error) {
	b := make([]byte, unsafe.Sizeof(RegPacketAck{}))
	headerBuff, err := packet.Encode(ack.header)
	if err != nil {
		return nil, err
	}
	var idx = 0
	idx = packet.EncodeBytes(b, headerBuff, idx)
	idx = packet.EncodeBytes(b, ack.RegMac, idx)
	idx = packet.EncodeBytes(b, ack.AutoIP, idx)
	idx = packet.EncodeBytes(b, ack.Mask, idx)
	return b, nil
}

func Decode(udpBytes []byte) (RegPacketAck, error) {
	size := unsafe.Sizeof(packet.Header{})
	res := RegPacketAck{}
	h, err := packet.Decode(udpBytes[:size])
	if err != nil {
		return RegPacketAck{}, err
	}
	var idx = 0
	res.header = h
	idx += int(size)
	mac := make([]byte, packet.MAC_SIZE)
	idx = packet.DecodeBytes(&mac, udpBytes, idx)
	res.RegMac = mac
	ip := make([]byte, packet.IP_SIZE)
	idx = packet.DecodeBytes(&ip, udpBytes, idx)
	res.AutoIP = ip
	mask := make([]byte, packet.IP_SIZE)
	idx = packet.DecodeBytes(&mask, udpBytes, idx)
	res.Mask = mask
	return res, nil
}
