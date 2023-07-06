package handshake

import (
	"errors"
	"github.com/topcloudz/fvpn/pkg/handler"
	"github.com/topcloudz/fvpn/pkg/packet"
	"github.com/topcloudz/fvpn/pkg/packet/header"
	"github.com/topcloudz/fvpn/pkg/util"
	"net"
	"unsafe"
)

type HandShakePacket struct {
	header header.Header //20
	SrcIP  net.IP
	PubKey [32]byte //dh public key, generate from curve25519
}

func NewPacket(msgType uint16, userId string) HandShakePacket {
	headerPacket, _ := header.NewHeader(msgType, userId)
	return HandShakePacket{
		header: headerPacket,
		SrcIP:  net.IP{},
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
	idx = packet.EncodeBytes(b, np.SrcIP[:], idx)
	idx = packet.EncodeBytes(b, np.PubKey[:], idx)

	return b, nil
}

func Decode(buff []byte) (HandShakePacket, error) {
	res := NewPacket(util.HandShakeMsgType, handler.UCTL.UserId)
	h, err := header.Decode(buff)
	if err != nil {
		return HandShakePacket{}, errors.New("decode common packet failed")
	}
	idx := 0
	res.header = h
	idx += int(unsafe.Sizeof(header.Header{}))

	srcIP := make([]byte, 16)
	idx = packet.DecodeBytes(&srcIP, buff, idx)
	copy(res.SrcIP[:], srcIP)

	pubKey := make([]byte, 32)
	idx = packet.DecodeBytes(&pubKey, buff, idx)
	copy(res.PubKey[:], pubKey[:])

	return res, nil
}
