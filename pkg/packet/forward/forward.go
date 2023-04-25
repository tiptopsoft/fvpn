package forward

import (
	"errors"
	"github.com/topcloudz/fvpn/pkg/option"
	"github.com/topcloudz/fvpn/pkg/packet"
	"net"
	"unsafe"
)

// ForwardPacket is through packet used in server
type ForwardPacket struct {
	header packet.Header
	body
}

type body struct {
	SrcMac net.HardwareAddr
	DstMac net.HardwareAddr
}

func NewPacket(networkId string) ForwardPacket {
	header := packet.NewHeader(option.MsgTypePacket, networkId)
	return ForwardPacket{
		header: header,
	}
}

func (fp ForwardPacket) Encode() ([]byte, error) {
	b := make([]byte, unsafe.Sizeof(ForwardPacket{}))
	headerBuff, err := fp.header.Encode()
	if err != nil {
		return nil, errors.New("encode common packet failed")
	}
	idx := 0
	idx = packet.EncodeBytes(b, headerBuff, idx)
	idx = packet.EncodeBytes(b, fp.SrcMac[:], idx)
	idx = packet.EncodeBytes(b, fp.DstMac[:], idx)
	return b, nil
}

func (fp ForwardPacket) Decode(udpBytes []byte) (packet.Interface, error) {
	res := NewPacket("")
	header, err := packet.NewPacketWithoutType().Decode(udpBytes)
	if err != nil {
		return ForwardPacket{}, errors.New("decode header packet failed")
	}
	idx := 0
	res.header = header.(packet.Header)
	idx += int(unsafe.Sizeof(packet.Header{}))
	var srcMac = make([]byte, 6)
	idx = packet.DecodeBytes(&srcMac, udpBytes, idx)
	res.SrcMac = srcMac
	var dstMac = make([]byte, 6)
	packet.DecodeBytes(&dstMac, udpBytes, idx)
	res.DstMac = dstMac
	return res, nil
}

func (fp ForwardPacket) Get() ForwardPacket {
	return fp
}
