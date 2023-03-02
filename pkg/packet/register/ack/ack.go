package ack

import (
	"github.com/interstellar-cloud/star/pkg/option"
	packet "github.com/interstellar-cloud/star/pkg/packet"
	"github.com/interstellar-cloud/star/pkg/packet/common"
	common2 "github.com/interstellar-cloud/star/pkg/packet/common"
	"net"
	"unsafe"
)

type RegPacketAck struct {
	common2.CommonPacket                  //8 byte
	RegMac               net.HardwareAddr //6 byte
	AutoIP               net.IP           //4byte
	Mask                 net.IP
}

func NewPacket() RegPacketAck {
	cmPacket := common2.NewPacket(option.MsgTypeRegisterAck)
	return RegPacketAck{
		CommonPacket: cmPacket,
	}
}

func (r RegPacketAck) Encode() ([]byte, error) {
	b := make([]byte, unsafe.Sizeof(RegPacketAck{}))
	cp, err := r.CommonPacket.Encode()
	if err != nil {
		return nil, err
	}
	var idx = 0
	idx = packet.EncodeBytes(b, cp, idx)
	idx = packet.EncodeBytes(b, r.RegMac, idx)
	idx = packet.EncodeBytes(b, r.AutoIP, idx)
	idx = packet.EncodeBytes(b, r.Mask, idx)
	return b, nil
}

func (r RegPacketAck) Decode(udpBytes []byte) (packet.Interface, error) {
	size := unsafe.Sizeof(common.CommonPacket{})
	res := RegPacketAck{}
	p, err := common.NewPacketWithoutType().Decode(udpBytes[:size])
	if err != nil {
		return RegPacketAck{}, err
	}
	var idx = 0
	res.CommonPacket = p.(common.CommonPacket)
	idx += int(size)
	mac := make([]byte, packet.MAC_SIZE)
	idx = packet.DecodeBytes(&mac, udpBytes, idx)
	res.RegMac = mac
	ip := make([]byte, packet.IP_SIZE)
	idx = packet.DecodeBytes(&ip, udpBytes, idx)
	res.AutoIP = ip
	mask := make([]byte, packet.IP_SIZE)
	idx = packet.DecodeBytes(&mask, udpBytes, idx)
	res.Mask = mask
	return res, nil
}
