package ack

import "unsafe"

type RegPacketAck struct {
	RegMac [4]byte
	AutoIP [4]byte
	Mask   [4]byte
}

func NewPacket() RegPacketAck {
	return RegPacketAck{}
}

func (regAck RegPacketAck) Encode(reg RegPacketAck) ([]byte, error) {
	b := make([]byte, unsafe.Sizeof(RegPacketAck{}))
	copy(b[20:24], reg.RegMac[:])
	return b, nil
}

func (regAck RegPacketAck) Decode(udpBytes []byte) (RegPacketAck, error) {
	res := RegPacketAck{}
	copy(regAck.RegMac[:], udpBytes[20:24])
	return res, nil
}
