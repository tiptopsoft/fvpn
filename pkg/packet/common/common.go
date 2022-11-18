package common

import (
	"github.com/interstellar-cloud/star/pkg/packet"
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

//CommonPacket  every time sends base frame.
type CommonPacket struct {
	Version uint8   //1
	TTL     uint8   //1
	Flags   uint16  //2
	Group   [4]byte //4
}

func NewPacket() CommonPacket {
	return CommonPacket{
		Version: Version,
		TTL:     DefaultTTL,
		Flags:   0,
		Group:   [4]byte{},
	}
}

func (cp CommonPacket) Encode() ([]byte, error) {

	idx := 0
	b := make([]byte, unsafe.Sizeof(cp))
	idx = packet.EncodeUint8(b, cp.Version, idx)
	idx = packet.EncodeUint8(b, cp.TTL, idx)
	idx = packet.EncodeUint16(b, cp.Flags, idx)
	packet.EncodeBytes(b, cp.Group[:], idx)
	return b, nil
}

func (cp CommonPacket) Decode(udpByte []byte) (CommonPacket, error) {
	idx := 0
	idx = packet.DecodeUint8(cp.Version, udpByte, idx)
	idx = packet.DecodeUint8(cp.TTL, udpByte, idx)
	idx = packet.DecodeUint16(cp.Flags, udpByte, idx)
	idx = packet.DecodeBytes(cp.Group[:], udpByte, idx)
	return cp, nil
}
