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

package peer

import (
	"bytes"
	"encoding/gob"
	"errors"
	"github.com/tiptopsoft/fvpn/pkg/packet"
	"github.com/tiptopsoft/fvpn/pkg/security"
	"github.com/tiptopsoft/fvpn/pkg/util"
	"net"
)

type PeerInfo struct {
	IP         net.IP
	NatIP      net.IP
	RemoteAddr net.UDPAddr
	PubKey     security.NoisePublicKey
}

type PeerPacket struct {
	Header packet.Header
	Peers  []PeerInfo
}

func NewPeerPacket(userId string) PeerPacket {
	h, _ := packet.NewHeader(util.MsgTypeQueryPeer, userId)
	return PeerPacket{
		Header: h,
		Peers:  nil,
	}
}

func Encode(peerPacket PeerPacket) ([]byte, error) {
	buff := make([]byte, packet.FvpnPktBuffSize)
	headerBuff, err := packet.Encode(peerPacket.Header)
	if err != nil {
		return nil, errors.New("encode common packet failed")
	}
	idx := 0
	idx = packet.EncodeBytes(buff, headerBuff, idx)
	buf := &bytes.Buffer{}
	b := gob.NewEncoder(buf)
	//err := binary.Write(buf, binary.BigEndian, peerPacket)
	err = b.Encode(peerPacket.Peers)
	if err != nil {
		return nil, err
	}

	idx = packet.EncodeBytes(buff, buf.Bytes(), idx)

	return buff[:idx], err
}

func Decode(buff []byte) (peerPacket PeerPacket, err error) {
	h, err := packet.Decode(buff)
	if err != nil {
		return PeerPacket{}, err
	}

	peerPacket = PeerPacket{}
	peerPacket.Header = h
	buf := bytes.NewReader(buff[packet.HeaderBuffSize:])
	//err = binary.Read(buf, binary.BigEndian, &peerPacket)
	d := gob.NewDecoder(buf)
	err = d.Decode(&peerPacket.Peers)
	if err != nil {
		return PeerPacket{}, err
	}

	return peerPacket, nil
}
