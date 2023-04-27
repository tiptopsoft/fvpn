package packet

import (
	"encoding/hex"
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

// Header  every time sends util frame. 12 byte
type Header struct {
	Version   uint8  //1
	TTL       uint8  //1
	Flags     uint16 //2
	NetworkId string //8
}

func NewPacketWithoutType() Header {
	return Header{
		Version: Version,
		TTL:     DefaultTTL,
		Flags:   0,
	}
}

func NewHeader(msgType uint16, networkId string) Header {
	return Header{
		Version:   Version,
		TTL:       DefaultTTL,
		Flags:     msgType,
		NetworkId: networkId,
	}
}

func (cp Header) Encode() ([]byte, error) {
	idx := 0
	b := make([]byte, unsafe.Sizeof(Header{}))
	idx = EncodeUint8(b, cp.Version, idx)
	idx = EncodeUint8(b, cp.TTL, idx)
	idx = EncodeUint16(b, cp.Flags, idx)
	buff, err := hex.DecodeString(cp.NetworkId)
	if err != nil {
		return nil, err
	}
	EncodeBytes(b, buff, idx)
	return b, nil
}

func (cp Header) Decode(udpByte []byte) (Interface, error) {
	idx := 0
	idx = DecodeUint8(&cp.Version, udpByte, idx)
	idx = DecodeUint8(&cp.TTL, udpByte, idx)
	idx = DecodeUint16(&cp.Flags, udpByte, idx)
	idx = DecodeNetworkId(cp.NetworkId, udpByte, idx)
	return cp, nil
}
