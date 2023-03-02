package register

import (
	"errors"
	"github.com/interstellar-cloud/star/pkg/option"
	packet "github.com/interstellar-cloud/star/pkg/packet"
	common "github.com/interstellar-cloud/star/pkg/packet/common"
	"net"
	"reflect"
	"unsafe"
)

// RegPacket registry a edge to registry
type RegPacket struct {
	header common.PacketHeader
	SrcMac net.HardwareAddr
}

func NewPacket() RegPacket {
	cmPacket := common.NewPacket(option.MsgTypeRegisterSuper)
	return RegPacket{
		header: cmPacket,
	}
}

func NewUnregisterPacket() RegPacket {
	cmPacket := common.NewPacket(option.MsgTypeUnregisterSuper)
	return RegPacket{
		header: cmPacket,
	}
}

func (r RegPacket) Encode() ([]byte, error) {
	b := make([]byte, unsafe.Sizeof(reflect.ValueOf(r)))
	commonBytes, err := r.header.Encode()
	if err != nil {
		return nil, errors.New("encode common packet failed")
	}
	idx := 0
	idx = packet.EncodeBytes(b, commonBytes, idx)
	idx = packet.EncodeBytes(b, r.SrcMac[:], idx)
	return b, nil
}

func (r RegPacket) Decode(buff []byte) (packet.Interface, error) {
	res := NewPacket()
	idx := 0
	idx += int(unsafe.Sizeof(common.PacketHeader{}))
	var mac = make([]byte, 6)
	packet.DecodeBytes(&mac, buff, idx)
	res.SrcMac = mac
	return res, nil
}
