package peer

import (
	"errors"
	"github.com/interstellar-cloud/star/pkg/packet"
	"github.com/interstellar-cloud/star/pkg/packet/common"
	"net"
	"reflect"
	"unsafe"
)

type PeerPacket struct {
	common.CommonPacket
	SrcMac net.HardwareAddr
}

func NewPacket() PeerPacket {
	return PeerPacket{}
}

func Encode(cp PeerPacket) ([]byte, error) {
	b := make([]byte, unsafe.Sizeof(reflect.ValueOf(cp)))
	commonBytes, err := common.Encode(cp.CommonPacket)
	if err != nil {
		return nil, errors.New("encode common packet failed")
	}
	idx := 0
	idx = packet.EncodeBytes(b, commonBytes, idx)
	idx = packet.EncodeBytes(b, cp.SrcMac[:], idx)
	return b, nil
}

func Decode(udpBytes []byte) (PeerPacket, error) {

	res := NewPacket()
	cp, err := common.Decode(udpBytes)
	if err != nil {
		return PeerPacket{}, errors.New("decode common packet failed")
	}
	idx := 0
	res.CommonPacket = cp
	idx += int(unsafe.Sizeof(common.NewPacket()))
	var mac = make([]byte, 6)
	packet.DecodeBytes(&mac, udpBytes, idx)
	res.SrcMac = mac
	return res, nil
}

func DecodeWithCommonPacket(udpBytes []byte, cp common.CommonPacket) (PeerPacket, error) {

	res := NewPacket()
	idx := 0
	res.CommonPacket = cp
	idx += int(unsafe.Sizeof(common.NewPacket()))
	var mac = make([]byte, 6)
	packet.DecodeBytes(&mac, udpBytes, idx)
	res.SrcMac = mac
	return res, nil
}
