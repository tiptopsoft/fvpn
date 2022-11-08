package common

import "github.com/interstellar-cloud/star/pkg/packet"

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

func NewPacket() *CommonPacket {
	return &CommonPacket{
		Version: Version,
		TTL:     DefaultTTL,
		Flags:   0,
		Group:   [4]byte{},
	}
}

func (cp *CommonPacket) Encode() ([]byte, error) {

	var b [8]byte
	b[0] = cp.Version
	copy(b[1:2], []byte{cp.TTL})
	copy(b[2:4], packet.EncodeUint16(cp.Flags))
	copy(b[4:8], cp.Group[:])
	return b[:], nil
}

func (cp *CommonPacket) Decode(udpByte []byte) (*CommonPacket, error) {
	cp.Version = udpByte[0]
	cp.TTL = udpByte[1]
	cp.Flags = packet.BytesToInt16(udpByte[2:4])
	copy(cp.Group[:], udpByte[4:8])

	return cp, nil
}
