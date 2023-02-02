package forward

import (
	"errors"
	"github.com/interstellar-cloud/star/pkg/packet"
	"github.com/interstellar-cloud/star/pkg/packet/common"
	"github.com/interstellar-cloud/star/pkg/util/option"
	"net"
	"unsafe"
)

//ForwardPacket is through packet used in registry
type ForwardPacket struct {
	common.CommonPacket
	SrcMac net.HardwareAddr
	DstMac net.HardwareAddr
}

func NewPacket() ForwardPacket {
	cmPacket := common.NewPacket(option.MsgTypePacket)
	return ForwardPacket{
		CommonPacket: cmPacket,
	}
}

func Encode(fp ForwardPacket) ([]byte, error) {
	b := make([]byte, unsafe.Sizeof(ForwardPacket{}))
	commonBytes, err := common.Encode(fp.CommonPacket)
	if err != nil {
		return nil, errors.New("encode common packet failed")
	}
	idx := 0
	idx = packet.EncodeBytes(b, commonBytes, idx)
	idx = packet.EncodeBytes(b, fp.SrcMac[:], idx)
	idx = packet.EncodeBytes(b, fp.DstMac[:], idx)
	return nil, nil
}

func Decode(udpBytes []byte) (ForwardPacket, error) {
	res := NewPacket()
	cp, err := common.Decode(udpBytes)
	if err != nil {
		return ForwardPacket{}, errors.New("decode common packet failed")
	}
	idx := 0
	res.CommonPacket = cp
	idx += int(unsafe.Sizeof(common.CommonPacket{}))
	var srcMac = make([]byte, 6)
	idx = packet.DecodeBytes(&srcMac, udpBytes, idx)
	res.SrcMac = srcMac
	var dstMac = make([]byte, 6)
	packet.DecodeBytes(&dstMac, udpBytes, idx)
	res.DstMac = dstMac
	return res, nil
}
