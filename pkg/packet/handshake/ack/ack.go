package ack

import (
	"errors"
	"github.com/topcloudz/fvpn/pkg/packet"
	"github.com/topcloudz/fvpn/pkg/packet/header"
	"github.com/topcloudz/fvpn/pkg/util"
	"unsafe"
)

type HandShakePacketAck struct {
	header header.Header //12
	PubKey [32]byte      //dh public key, generate from curve25519
	//SrcIP  net.IP
	//PubKey  [16]byte
}

func NewPacket() HandShakePacketAck {
	headerPacket, _ := header.NewHeader(util.HandShakeMsgType, "")
	return HandShakePacketAck{
		header: headerPacket,
	}
}

func Encode(np HandShakePacketAck) ([]byte, error) {
	b := make([]byte, unsafe.Sizeof(HandShakePacketAck{}))
	headerBuff, err := header.Encode(np.header)
	if err != nil {
		return nil, errors.New("encode common packet failed")
	}
	idx := 0
	idx = packet.EncodeBytes(b, headerBuff, idx)
	//idx = packet.EncodeBytes(b, np.SrcIP[:], idx)
	//idx = packet.EncodeBytes(b, np.PubKey[:], idx)
	idx = packet.EncodeBytes(b, np.PubKey[:], idx)

	return b, nil
}

func Decode(buff []byte) (HandShakePacketAck, error) {
	res := NewPacket()
	h, err := header.Decode(buff)
	if err != nil {
		return HandShakePacketAck{}, errors.New("decode common packet failed")
	}
	idx := 0
	res.header = h
	idx += int(unsafe.Sizeof(header.Header{}))

	//srcIP := make([]byte, 16)
	//idx = packet.DecodeBytes(&srcIP, buff, idx)
	//copy(res.SrcIP, srcIP)
	//
	//appId := make([]byte, 16)
	//idx = packet.DecodeBytes(&appId, buff, idx)
	//copy(res.PubKey[:], appId)

	pubKey := make([]byte, 32)
	idx = packet.DecodeBytes(&pubKey, buff, idx)
	copy(res.PubKey[:], pubKey[:])

	return res, nil
}
