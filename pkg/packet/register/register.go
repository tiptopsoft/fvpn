package register

import (
	"errors"
	"github.com/interstellar-cloud/star/pkg/option"
	"github.com/interstellar-cloud/star/pkg/packet"
	"github.com/interstellar-cloud/star/pkg/packet/common"
	"net"
	"reflect"
	"unsafe"
)

// RegPacket registry a edge to registry
type RegPacket struct {
	common.CommonPacket
	SrcMac net.HardwareAddr
}

func NewPacket() RegPacket {
	cmPacket := common.NewPacket(option.MsgTypeRegisterSuper)
	return RegPacket{
		CommonPacket: cmPacket,
	}
}

func NewUnregisterPacket() RegPacket {
	cmPacket := common.NewPacket(option.MsgTypeUnregisterSuper)
	return RegPacket{
		CommonPacket: cmPacket,
	}
}

func Encode(cp RegPacket) ([]byte, error) {
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

func Decode(udpBytes []byte) (RegPacket, error) {

	res := NewPacket()
	cp, err := common.Decode(udpBytes)
	if err != nil {
		return RegPacket{}, errors.New("decode common packet failed")
	}
	idx := 0
	res.CommonPacket = cp
	idx += int(unsafe.Sizeof(common.CommonPacket{}))
	var mac = make([]byte, 6)
	packet.DecodeBytes(&mac, udpBytes, idx)
	res.SrcMac = mac
	return res, nil
}

func DecodeWithCommonPacket(udpBytes []byte, cp common.CommonPacket) (RegPacket, error) {

	res := NewPacket()
	idx := 0
	res.CommonPacket = cp
	idx += int(unsafe.Sizeof(common.CommonPacket{}))
	var mac = make([]byte, 6)
	packet.DecodeBytes(&mac, udpBytes, idx)
	res.SrcMac = mac
	return res, nil
}
