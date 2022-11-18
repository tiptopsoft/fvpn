package register

import (
	"errors"
	"github.com/interstellar-cloud/star/pkg/packet"
	"github.com/interstellar-cloud/star/pkg/packet/common"
	"net"
	"reflect"
	"unsafe"
)

// RegPacket register a edge to register
type RegPacket struct {
	common.CommonPacket
	SrcMac net.HardwareAddr
}

func NewPacket() RegPacket {
	return RegPacket{}
}

func (cp RegPacket) Encode() ([]byte, error) {
	b := make([]byte, unsafe.Sizeof(reflect.ValueOf(cp)))
	commonBytes, err := cp.CommonPacket.Encode()
	if err != nil {
		return nil, errors.New("encode common packet failed")
	}
	idx := 0
	idx = packet.EncodeBytes(b, commonBytes, idx)
	idx = packet.EncodeBytes(b, cp.SrcMac[:], idx)
	return b, nil
}

func (reg RegPacket) Decode(udpBytes []byte) (RegPacket, error) {

	res := RegPacket{}
	cp, err := common.NewPacket().Decode(udpBytes)
	if err != nil {
		return RegPacket{}, errors.New("decode common packet failed")
	}
	idx := 0
	res.CommonPacket = cp
	idx += int(unsafe.Sizeof(common.NewPacket()))
	var mac = make([]byte, 6)
	packet.DecodeBytes(&mac, udpBytes, idx)
	res.SrcMac = mac
	return res, nil
}
