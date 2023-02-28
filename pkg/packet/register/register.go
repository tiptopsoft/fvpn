package register

import (
	"errors"
	"github.com/interstellar-cloud/star/pkg/option"
	packet2 "github.com/interstellar-cloud/star/pkg/packet"
	common2 "github.com/interstellar-cloud/star/pkg/packet/common"
	"net"
	"reflect"
	"unsafe"
)

// RegPacket registry a edge to registry
type RegPacket struct {
	common2.CommonPacket
	SrcMac net.HardwareAddr
}

func NewPacket() RegPacket {
	cmPacket := common2.NewPacket(option.MsgTypeRegisterSuper)
	return RegPacket{
		CommonPacket: cmPacket,
	}
}

func NewUnregisterPacket() RegPacket {
	cmPacket := common2.NewPacket(option.MsgTypeUnregisterSuper)
	return RegPacket{
		CommonPacket: cmPacket,
	}
}

func (r RegPacket) Encode() ([]byte, error) {
	b := make([]byte, unsafe.Sizeof(reflect.ValueOf(r)))
	commonBytes, err := r.CommonPacket.Encode()
	if err != nil {
		return nil, errors.New("encode common packet failed")
	}
	idx := 0
	idx = packet2.EncodeBytes(b, commonBytes, idx)
	idx = packet2.EncodeBytes(b, r.SrcMac[:], idx)
	return b, nil
}

func (r RegPacket) Decode(buff []byte) (packet2.Interface, error) {
	res := NewPacket()
	idx := 0
	idx += int(unsafe.Sizeof(common2.CommonPacket{}))
	var mac = make([]byte, 6)
	packet2.DecodeBytes(&mac, buff, idx)
	res.SrcMac = mac
	return res, nil
}
