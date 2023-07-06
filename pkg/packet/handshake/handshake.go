package handshake

import (
	"errors"
	"github.com/topcloudz/fvpn/pkg/handler"
	"github.com/topcloudz/fvpn/pkg/packet"
	"github.com/topcloudz/fvpn/pkg/util"
)

type HandShakePacket struct {
	Header packet.Header
	PubKey [32]byte //dh public key, generate from curve25519
}

func NewPacket(msgType uint16, userId string) HandShakePacket {
	headerPacket, _ := packet.NewHeader(msgType, userId)
	return HandShakePacket{
		Header: headerPacket,
	}
}

func Encode(np HandShakePacket) ([]byte, error) {
	b := make([]byte, packet.HandshakeBuffSize)
	headerBuff, err := packet.Encode(np.Header)
	if err != nil {
		return nil, errors.New("encode common packet failed")
	}
	idx := 0
	idx = packet.EncodeBytes(b, headerBuff, idx)
	idx = packet.EncodeBytes(b, np.PubKey[:], idx)

	return b, nil
}

func Decode(buff []byte) (HandShakePacket, error) {
	res := NewPacket(util.HandShakeMsgType, handler.UCTL.UserId)
	h, err := packet.Decode(buff)
	if err != nil {
		return HandShakePacket{}, errors.New("decode common packet failed")
	}
	idx := 0
	res.Header = h
	idx += packet.HeaderBuffSize

	pubKey := make([]byte, 32)
	idx = packet.DecodeBytes(&pubKey, buff, idx)
	copy(res.PubKey[:], pubKey[:])

	return res, nil
}
