package common

import (
	"github.com/interstellar-cloud/star/pkg/util/packet"
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

func NewPacket(msgType uint16) CommonPacket {
	return CommonPacket{
		Version: Version,
		TTL:     DefaultTTL,
		Flags:   msgType,
		Group:   [4]byte{},
	}
}

func Encode(cp CommonPacket) ([]byte, error) {

	idx := 0
	b := make([]byte, unsafe.Sizeof(CommonPacket{}))
	idx = packet.EncodeUint8(b, cp.Version, idx)
	idx = packet.EncodeUint8(b, cp.TTL, idx)
	idx = packet.EncodeUint16(b, cp.Flags, idx)
	packet.EncodeBytes(b, cp.Group[:], idx)
	return b, nil
}

func Decode(udpByte []byte) (CommonPacket, error) {
	cp := CommonPacket{}
	idx := 0
	idx = packet.DecodeUint8(&cp.Version, udpByte, idx)
	idx = packet.DecodeUint8(&cp.TTL, udpByte, idx)
	idx = packet.DecodeUint16(&cp.Flags, udpByte, idx)
	a := cp.Group[:]
	idx = packet.DecodeBytes(&a, udpByte, idx)
	return cp, nil
}
