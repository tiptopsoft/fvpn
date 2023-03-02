package forward

import (
	"errors"
	"github.com/interstellar-cloud/star/pkg/option"
	packet2 "github.com/interstellar-cloud/star/pkg/packet"
	"github.com/interstellar-cloud/star/pkg/packet/common"
	common2 "github.com/interstellar-cloud/star/pkg/packet/common"
	"net"
	"unsafe"
)

//ForwardPacket is through packet used in registry
type ForwardPacket struct {
	common2.CommonPacket
	SrcMac net.HardwareAddr
	DstMac net.HardwareAddr
}

func NewPacket() ForwardPacket {
	cmPacket := common2.NewPacket(option.MsgTypePacket)
	return ForwardPacket{
		CommonPacket: cmPacket,
	}
}

func (fp ForwardPacket) Encode() ([]byte, error) {
	b := make([]byte, unsafe.Sizeof(ForwardPacket{}))
	commonBytes, err := fp.CommonPacket.Encode()
	if err != nil {
		return nil, errors.New("encode common packet failed")
	}
	idx := 0
	idx = packet2.EncodeBytes(b, commonBytes, idx)
	idx = packet2.EncodeBytes(b, fp.SrcMac[:], idx)
	idx = packet2.EncodeBytes(b, fp.DstMac[:], idx)
	return b, nil
}

func (fp ForwardPacket) Decode(udpBytes []byte) (packet2.Interface, error) {
	res := NewPacket()
	cp, err := common.NewPacketWithoutType().Decode(udpBytes)
	if err != nil {
		return ForwardPacket{}, errors.New("decode common packet failed")
	}
	idx := 0
	res.CommonPacket = cp.(common.CommonPacket)
	idx += int(unsafe.Sizeof(common2.CommonPacket{}))
	var srcMac = make([]byte, 6)
	idx = packet2.DecodeBytes(&srcMac, udpBytes, idx)
	res.SrcMac = srcMac
	var dstMac = make([]byte, 6)
	packet2.DecodeBytes(&dstMac, udpBytes, idx)
	res.DstMac = dstMac
	return res, nil
}

func (fp ForwardPacket) Get() ForwardPacket {
	return fp
}
