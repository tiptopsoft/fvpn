package packet

import (
	"encoding/hex"
	"errors"
	"net"
	"unsafe"
)

const (
	Version    uint8 = 1
	DefaultTTL uint8 = 100
)

// Header  every time sends util frame. 44 byte
type Header struct {
	Version uint8  //1
	TTL     uint8  //1
	Flags   uint16 //2
	UserId  [8]byte
	SrcIP   net.IP //16
	DstIP   net.IP //16
}

func NewHeader(msgType uint16, userId string) (Header, error) {
	bs, err := hex.DecodeString(userId)
	if err != nil {
		return Header{}, errors.New("invalid userId")
	}

	h := Header{
		Version: Version,
		TTL:     DefaultTTL,
		Flags:   msgType,
	}
	copy(h.UserId[:], bs)
	return h, nil
}

func Encode(h Header) ([]byte, error) {
	idx := 0
	b := make([]byte, unsafe.Sizeof(Header{}))
	idx = EncodeUint8(b, h.Version, idx)
	idx = EncodeUint8(b, h.TTL, idx)
	idx = EncodeUint16(b, h.Flags, idx)
	idx = EncodeBytes(b, h.UserId[:], idx)
	idx = EncodeBytes(b, h.SrcIP[:], idx)
	idx = EncodeBytes(b, h.DstIP[:], idx)
	return b, nil
}

func Decode(buff []byte) (h Header, err error) {
	idx := 0
	idx = DecodeUint8(&h.Version, buff, idx)
	idx = DecodeUint8(&h.TTL, buff, idx)
	idx = DecodeUint16(&h.Flags, buff, idx)
	userId := make([]byte, 8)
	idx = DecodeBytes(&userId, buff, idx)
	copy(h.UserId[:], userId)

	srcIP := make([]byte, 16)
	idx = DecodeBytes(&srcIP, buff, idx)
	h.SrcIP = srcIP

	dstIP := make([]byte, 16)
	idx = DecodeBytes(&dstIP, buff, idx)
	h.DstIP = dstIP
	return
}
