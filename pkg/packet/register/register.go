package register

import (
	"errors"
	"github.com/topcloudz/fvpn/pkg/packet"
	"github.com/topcloudz/fvpn/pkg/packet/header"
	"github.com/topcloudz/fvpn/pkg/util"
	"net"
	"unsafe"
)

// RegPacket server a client to server
type RegPacket struct { //48
	header header.Header //12
	SrcIP  net.IP
	PubKey [16]byte
	UserId [10]byte
}

func NewPacket() RegPacket {
	cmPacket, _ := header.NewHeader(util.MsgTypeRegisterSuper, "")
	reg := RegPacket{
		header: cmPacket,
	}
	return reg
}

func Encode(regPacket RegPacket) ([]byte, error) {
	b := make([]byte, 48)
	headerBuff, err := header.Encode(regPacket.header)
	if err != nil {
		return nil, errors.New("encode Header failed")
	}
	idx := 0
	idx = packet.EncodeBytes(b, headerBuff, idx)
	idx = packet.EncodeBytes(b, regPacket.PubKey[:], idx)
	idx = packet.EncodeBytes(b, regPacket.SrcIP, idx)
	return b, nil
}

func Decode(buff []byte) (RegPacket, error) {
	res := NewPacket()

	h, err := header.Decode(buff[:12])
	if err != nil {
		return RegPacket{}, err
	}
	res.header = h
	idx := 0
	idx += int(unsafe.Sizeof(header.Header{}))
	var appId = make([]byte, 16)
	idx = packet.DecodeBytes(&appId, buff, idx)
	copy(res.PubKey[:], appId)
	var ip = make([]byte, 16)
	idx = packet.DecodeBytes(&ip, buff, idx)
	copy(res.SrcIP[:], ip)
	return res, nil
}
