package ack

import (
	"github.com/interstellar-cloud/star/pkg/packet"
	"github.com/interstellar-cloud/star/pkg/packet/common"
	"net"
	"unsafe"
)

type RegPacketAck struct {
	common.CommonPacket                  //8 byte
	RegMac              net.HardwareAddr //6 byte
	AutoIP              net.IP           //4byte
	Mask                net.IPMask
}

func NewPacket() RegPacketAck {
	return RegPacketAck{}
}

func Encode(reg RegPacketAck) ([]byte, error) {
	b := make([]byte, unsafe.Sizeof(RegPacketAck{}))
	cp, err := common.Encode(reg.CommonPacket)
	if err != nil {
		return nil, err
	}
	var idx = 0
	idx = packet.EncodeBytes(b, cp, idx)
	idx = packet.EncodeBytes(b, reg.RegMac, idx)
	idx = packet.EncodeBytes(b, reg.AutoIP, idx)
	idx = packet.EncodeBytes(b, reg.Mask, idx)
	return b, nil
}

func Decode(udpBytes []byte) (RegPacketAck, error) {
	res := RegPacketAck{}
	p, err := common.Decode(udpBytes[:8])
	if err != nil {
		return RegPacketAck{}, err
	}
	var idx = 0
	res.CommonPacket = p
	idx += int(unsafe.Sizeof(p))
	mac := make([]byte, 6)
	idx = packet.DecodeBytes(&mac, udpBytes, idx)
	res.RegMac = mac
	ip := make([]byte, 16)
	idx = packet.DecodeBytes(&ip, udpBytes, idx)
	res.AutoIP = ip
	mask := make([]byte, 16)
	idx = packet.DecodeBytes(&mask, udpBytes, idx)
	res.Mask = mask
	return res, nil
}
