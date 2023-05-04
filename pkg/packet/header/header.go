package header

import (
	"encoding/hex"
	"errors"
	"github.com/topcloudz/fvpn/pkg/packet"
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
	Version   uint8   //1
	TTL       uint8   //1
	Flags     uint16  //2
	NetworkId [8]byte //8
}

func NewPacketWithoutType() *Header {
	return &Header{
		Version: Version,
		TTL:     DefaultTTL,
		Flags:   0,
	}
}

func NewHeader(msgType uint16, networkId string) (Header, error) {
	bs, err := hex.DecodeString(networkId)
	if err != nil {
		return Header{}, errors.New("invalid networkId")
	}
	h := Header{
		Version: Version,
		TTL:     DefaultTTL,
		Flags:   msgType,
	}
	copy(h.NetworkId[:], bs)
	return h, nil
}

func Encode(h Header) ([]byte, error) {
	idx := 0
	b := make([]byte, unsafe.Sizeof(Header{}))
	idx = packet.EncodeUint8(b, h.Version, idx)
	idx = packet.EncodeUint8(b, h.TTL, idx)
	idx = packet.EncodeUint16(b, h.Flags, idx)
	idx = packet.EncodeBytes(b, h.NetworkId[:], idx)
	return b, nil
}

func Decode(udpByte []byte) (h Header, err error) {
	idx := 0
	idx = packet.DecodeUint8(&h.Version, udpByte, idx)
	idx = packet.DecodeUint8(&h.TTL, udpByte, idx)
	idx = packet.DecodeUint16(&h.Flags, udpByte, idx)
	b := make([]byte, 8)
	idx = packet.DecodeBytes(&b, udpByte, idx)
	copy(h.NetworkId[:], b)
	return
}
