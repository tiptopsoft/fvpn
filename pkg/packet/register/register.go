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

package register

import (
	"errors"
	"github.com/tiptopsoft/fvpn/pkg/packet"
	"github.com/tiptopsoft/fvpn/pkg/util"
	"net"
	"unsafe"
)

// RegPacket server a client to server
type RegPacket struct { //48
	header packet.Header //12
	SrcIP  net.IP
	PubKey [16]byte
	UserId [8]byte
}

func NewPacket() RegPacket {
	cmPacket, _ := packet.NewHeader(util.MsgTypeRegisterSuper, "")
	reg := RegPacket{
		header: cmPacket,
	}
	return reg
}

func Encode(regPacket RegPacket) ([]byte, error) {
	b := make([]byte, 48)
	headerBuff, err := packet.Encode(regPacket.header)
	if err != nil {
		return nil, errors.New("encode Header failed")
	}
	idx := 0
	idx = packet.EncodeBytes(b, headerBuff, idx)
	idx = packet.EncodeBytes(b, regPacket.PubKey[:], idx)
	idx = packet.EncodeBytes(b, regPacket.SrcIP, idx)
	return b, nil
}

func Decode(buff []byte) (RegPacket, error) {
	res := NewPacket()

	h, err := packet.Decode(buff[:packet.HeaderBuffSize])
	if err != nil {
		return RegPacket{}, err
	}
	res.header = h
	idx := 0
	idx += int(unsafe.Sizeof(packet.Header{}))
	var appId = make([]byte, 16)
	idx = packet.DecodeBytes(&appId, buff, idx)
	copy(res.PubKey[:], appId)
	var ip = make([]byte, 16)
	idx = packet.DecodeBytes(&ip, buff, idx)
	copy(res.SrcIP[:], ip)
	return res, nil
}
