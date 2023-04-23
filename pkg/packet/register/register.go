package register

import (
	"errors"
	"github.com/topcloudz/fvpn/pkg/option"
	packet "github.com/topcloudz/fvpn/pkg/packet"
	"net"
	"reflect"
	"unsafe"
)

// RegPacket server a client to server
type RegPacket struct {
	header packet.Header
	SrcMac net.HardwareAddr
}

func NewPacket() RegPacket {
	cmPacket := packet.NewHeader(option.MsgTypeRegisterSuper)
	return RegPacket{
		header: cmPacket,
	}
}

func NewUnregisterPacket() RegPacket {
	cmPacket := packet.NewHeader(option.MsgTypeUnregisterSuper)
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
	idx += int(unsafe.Sizeof(packet.Header{}))
	var mac = make([]byte, 6)
	packet.DecodeBytes(&mac, buff, idx)
	res.SrcMac = mac
	return res, nil
}
