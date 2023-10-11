// Copyright 2023 TiptopSoft, Inc.
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

package handshake

import (
	"errors"
	"github.com/tiptopsoft/fvpn/pkg/packet"
	"github.com/tiptopsoft/fvpn/pkg/util"
)

type Packet struct {
	Header packet.Header
	PubKey [32]byte //dh public key, generate from curve25519
}

func NewPacket(msgType uint16, userId string) Packet {
	headerPacket, _ := packet.NewHeader(msgType, userId)
	return Packet{
		Header: headerPacket,
	}
}

func (np Packet) Encode(buff []byte) (int, error) {
	b := buff[:packet.HandshakeBuffSize]
	idx, err := np.Header.Encode(b)
	if err != nil {
		return 0, errors.New("encode header failed")
	}
	idx = packet.EncodeBytes(b, np.PubKey[:], idx)

	return idx, nil
}

func Decode(msgType uint16, buff []byte) (Packet, error) {
	res := NewPacket(msgType, util.Info().GetUserId())
	h, err := packet.Decode(buff)
	if err != nil {
		return Packet{}, errors.New("decode header failed")
	}
	idx := 0
	res.Header = h
	idx += packet.HeaderBuffSize

	pubKey := make([]byte, 32)
	idx = packet.DecodeBytes(&pubKey, buff, idx)
	copy(res.PubKey[:], pubKey[:])

	return res, nil
}
