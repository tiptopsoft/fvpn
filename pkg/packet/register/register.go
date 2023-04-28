package register

import (
	"errors"
	"github.com/topcloudz/fvpn/pkg/option"
	packet "github.com/topcloudz/fvpn/pkg/packet"
	"net"
	"unsafe"
)

// RegPacket server a client to server
type RegPacket struct { //48
	header *packet.Header   //12
	SrcMac net.HardwareAddr //20
	SrcIP  net.IP           // 4 byte是ipv4, 16 byte是ipv6
}

func NewPacket(networkId string, srcMac net.HardwareAddr, srcIP net.IP) RegPacket {
	cmPacket, _ := packet.NewHeader(option.MsgTypeRegisterSuper, networkId)
	reg := RegPacket{
		header: cmPacket,
		SrcIP:  srcIP,
		SrcMac: srcMac,
	}

	return reg
}

func NewUnregisterPacket(networkId string) RegPacket {
	cmPacket, _ := packet.NewHeader(option.MsgTypeUnregisterSuper, networkId)
	return RegPacket{
		header: cmPacket,
	}
}

func (r RegPacket) Encode() ([]byte, error) {
	b := make([]byte, 48)
	commonBytes, err := r.header.Encode()
	if err != nil {
		return nil, errors.New("encode common packet failed")
	}
	idx := 0
	idx = packet.EncodeBytes(b, commonBytes, idx)
	idx = packet.EncodeBytes(b, r.SrcMac[:], idx)
	idx = packet.EncodeBytes(b, r.SrcIP[:], idx)
	return b, nil
}

func (r RegPacket) Decode(buff []byte) (packet.Interface, error) {
	res := NewPacket("", net.HardwareAddr{}, net.IP{})
	idx := 0
	idx += int(unsafe.Sizeof(packet.Header{}))
	var mac = make([]byte, 20)
	packet.DecodeBytes(&mac, buff, idx)
	copy(res.SrcMac[:], mac)
	var ip = make([]byte, 16)
	packet.DecodeBytes(&ip, buff, idx)
	copy(res.SrcIP[:], ip)
	return res, nil
}
