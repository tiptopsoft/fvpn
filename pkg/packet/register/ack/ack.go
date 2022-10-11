// Copyright 2023 Tiptopsoft, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package ack

import (
	"github.com/tiptopsoft/fvpn/pkg/packet"
	"github.com/tiptopsoft/fvpn/pkg/util"
	"net"
	"unsafe"
)

type RegPacketAck struct {
	header packet.Header    //8 byte
	RegMac net.HardwareAddr //6 byte
	AutoIP net.IP           //4byte
	Mask   net.IP
}

func NewPacket() RegPacketAck {
	cmPacket, _ := packet.NewHeader(util.MsgTypeRegisterAck, "")
	return RegPacketAck{
		header: cmPacket,
	}
}

func Encode(ack RegPacketAck) ([]byte, error) {
	b := make([]byte, unsafe.Sizeof(RegPacketAck{}))
	headerBuff, err := packet.Encode(ack.header)
	if err != nil {
		return nil, err
	}
	var idx = 0
	idx = packet.EncodeBytes(b, headerBuff, idx)
	idx = packet.EncodeBytes(b, ack.RegMac, idx)
	idx = packet.EncodeBytes(b, ack.AutoIP, idx)
	idx = packet.EncodeBytes(b, ack.Mask, idx)
	return b, nil
}

func Decode(udpBytes []byte) (RegPacketAck, error) {
	size := unsafe.Sizeof(packet.Header{})
	res := RegPacketAck{}
	h, err := packet.Decode(udpBytes[:size])
	if err != nil {
		return RegPacketAck{}, err
	}
	var idx = 0
	res.header = h
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
