package ack

import (
	"github.com/interstellar-cloud/star/pkg/packet/common"
	"unsafe"
)

type RegPacketAck struct {
	common.CommonPacket
	RegMac [4]byte
	AutoIP [4]byte
	Mask   [4]byte
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
	copy(b[:24], cp)
	copy(b[20:24], reg.RegMac[:])
	copy(b[24:28], reg.AutoIP[:])
	copy(b[28:32], reg.Mask[:])
	return b, nil
}

func (regAck RegPacketAck) Decode(udpBytes []byte) (RegPacketAck, error) {
	res := RegPacketAck{}
	p, err := common.NewPacket().Decode(udpBytes[:20])
	if err != nil {
		return RegPacketAck{}, err
	}
	regAck.CommonPacket = p
	copy(regAck.RegMac[:], udpBytes[20:24])
	copy(regAck.AutoIP[:], udpBytes[24:28])
	copy(regAck.Mask[:], udpBytes[28:32])
	return res, nil
}
