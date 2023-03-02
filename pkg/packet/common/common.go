package common

import (
	packet2 "github.com/interstellar-cloud/star/pkg/packet"
	"unsafe"
)

var (
	Version          uint8  = 1
	DefaultTTL       uint8  = 100
	IPV4             uint16 = 0x01
	IPV6             uint16 = 0x02
	COMMON_FRAM_SIZE        = 20
	DefaultPort      uint16 = 3000
)

//CommonPacket  every time sends util frame.
type CommonPacket struct {
	Version uint8   //1
	TTL     uint8   //1
	Flags   uint16  //2
	Group   [4]byte //4
}

func NewPacketWithoutType() CommonPacket {
	return CommonPacket{
		Version: Version,
		TTL:     DefaultTTL,
		Flags:   0,
		Group:   [4]byte{},
	}
}

func NewPacket(msgType uint16) CommonPacket {
	return CommonPacket{
		Version: Version,
		TTL:     DefaultTTL,
		Flags:   msgType,
		Group:   [4]byte{},
	}
}

func (cp CommonPacket) Encode() ([]byte, error) {
	idx := 0
	b := make([]byte, unsafe.Sizeof(CommonPacket{}))
	idx = packet2.EncodeUint8(b, cp.Version, idx)
	idx = packet2.EncodeUint8(b, cp.TTL, idx)
	idx = packet2.EncodeUint16(b, cp.Flags, idx)
	packet2.EncodeBytes(b, cp.Group[:], idx)
	return b, nil
}

func (cp CommonPacket) Decode(udpByte []byte) (packet2.Interface, error) {
	idx := 0
	idx = packet2.DecodeUint8(&cp.Version, udpByte, idx)
	idx = packet2.DecodeUint8(&cp.TTL, udpByte, idx)
	idx = packet2.DecodeUint16(&cp.Flags, udpByte, idx)
	a := cp.Group[:]
	idx = packet2.DecodeBytes(&a, udpByte, idx)
	return cp, nil
}
