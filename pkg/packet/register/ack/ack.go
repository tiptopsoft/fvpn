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
	Mask                net.IP           //4byte
}

func NewPacket() RegPacketAck {
	return RegPacketAck{}
}

func (regAck RegPacketAck) Encode(reg RegPacketAck) ([]byte, error) {
	b := make([]byte, unsafe.Sizeof(RegPacketAck{}))
	cp, err := reg.CommonPacket.Encode()
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

func (regAck RegPacketAck) Decode(udpBytes []byte) (RegPacketAck, error) {
	res := RegPacketAck{}
	p, err := common.NewPacket().Decode(udpBytes[:8])
	if err != nil {
		return RegPacketAck{}, err
	}
	var idx = 0
	res.CommonPacket = p
	idx += int(unsafe.Sizeof(p))
	idx = packet.DecodeBytes(&udpBytes, res.RegMac, idx)
	idx = packet.DecodeBytes(&udpBytes, res.AutoIP, idx)
	idx = packet.DecodeBytes(&udpBytes, res.Mask, idx)
	return res, nil
}
