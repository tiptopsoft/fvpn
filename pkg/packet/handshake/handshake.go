package handshake

import (
	"errors"
	"github.com/topcloudz/fvpn/pkg/option"
	"github.com/topcloudz/fvpn/pkg/packet"
	"github.com/topcloudz/fvpn/pkg/packet/header"
	"unsafe"
)

type HandShakePacket struct {
	header header.Header //12
}

func NewPacket(networkId string) HandShakePacket {
	headerPacket, _ := header.NewHeader(option.HandShakeMsgType, networkId)
	return HandShakePacket{
		header: headerPacket,
	}
}

func Encode(np HandShakePacket) ([]byte, error) {
	b := make([]byte, unsafe.Sizeof(HandShakePacket{}))
	headerBuff, err := header.Encode(np.header)
	if err != nil {
		return nil, errors.New("encode common packet failed")
	}
	idx := 0
	idx = packet.EncodeBytes(b, headerBuff, idx)

	return b, nil
}

func Decode(buff []byte) (HandShakePacket, error) {
	res := NewPacket("")
	h, err := header.Decode(buff)
	if err != nil {
		return HandShakePacket{}, errors.New("decode common packet failed")
	}
	idx := 0
	res.header = h
	idx += int(unsafe.Sizeof(header.Header{}))

	return res, nil
}