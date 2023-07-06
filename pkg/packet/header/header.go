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

// Header  every time sends util frame. 20 byte
type Header struct {
	Version   uint8   //1
	TTL       uint8   //1
	Flags     uint16  //2
	NetworkId [8]byte //8
	UserId    [8]byte
}

func NewHeader(msgType uint16, userId string) (Header, error) {
	bs, err := hex.DecodeString(userId)
	if err != nil {
		return Header{}, errors.New("invalid networkId")
	}

	h := Header{
		Version: Version,
		TTL:     DefaultTTL,
		Flags:   msgType,
	}
	copy(h.UserId[:], bs)
	//copy(h.PubKey[:], appIdData)
	return h, nil
}

func Encode(h Header) ([]byte, error) {
	idx := 0
	b := make([]byte, unsafe.Sizeof(Header{}))
	idx = packet.EncodeUint8(b, h.Version, idx)
	idx = packet.EncodeUint8(b, h.TTL, idx)
	idx = packet.EncodeUint16(b, h.Flags, idx)
	idx = packet.EncodeBytes(b, h.NetworkId[:], idx)
	idx = packet.EncodeBytes(b, h.UserId[:], idx)
	//idx = packet.EncodeBytes(b, h.PubKey[:], idx)
	return b, nil
}

func Decode(buff []byte) (h Header, err error) {
	idx := 0
	idx = packet.DecodeUint8(&h.Version, buff, idx)
	idx = packet.DecodeUint8(&h.TTL, buff, idx)
	idx = packet.DecodeUint16(&h.Flags, buff, idx)
	b := make([]byte, 8)
	idx = packet.DecodeBytes(&b, buff, idx)
	copy(h.NetworkId[:], b)
	u := make([]byte, 8)
	idx = packet.DecodeBytes(&u, buff, idx)
	copy(h.UserId[:], u)
	return
}
