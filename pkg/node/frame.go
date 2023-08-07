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

package node

import (
	"context"
	"encoding/hex"
	"github.com/tiptopsoft/fvpn/pkg/packet"
	"net"
	"sync"
)

type Frame struct {
	Ctx context.Context
	sync.Mutex
	Buff       []byte
	Packet     []byte
	Size       int
	NetworkId  string
	UserId     [8]byte
	RemoteAddr *net.UDPAddr //remote socket endpoint
	SrcIP      net.IP
	DstIP      net.IP
	FrameType  uint16
	Peer       *Peer
}

func (f *Frame) GetPeer() *Peer {
	return f.Peer
}

func NewFrame() *Frame {
	return &Frame{
		Ctx:    context.Background(),
		Buff:   make([]byte, packet.FvpnPktBuffSize),
		Packet: make([]byte, packet.FvpnPktBuffSize),
	}
}

func (f *Frame) Clear() {
	buf := make([]byte, packet.FvpnPktBuffSize)
	copy(f.Packet, buf)
}

func (f *Frame) UidString() string {
	return hex.EncodeToString(f.UserId[:])
}

func (f *Frame) Context() context.Context {
	return f.Ctx
}
