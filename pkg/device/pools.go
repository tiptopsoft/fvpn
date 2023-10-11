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

package device

import (
	"github.com/tiptopsoft/fvpn/pkg/packet"
	"golang.org/x/crypto/chacha20poly1305"
	"sync"
	"time"
)

type MemoryPool struct {
	lock sync.Mutex
	pool sync.Pool
	//cond sync.Cond
}

func InitPools() (buffPool *MemoryPool, framePool *MemoryPool) {
	framePool = NewPool(func() any {
		frame := &Frame{
			Ctx:        nil,
			Mutex:      sync.Mutex{},
			Packet:     make([]byte, packet.FvpnPktBuffSize),
			Buff:       make([]byte, packet.FvpnPktBuffSize),
			Size:       0,
			NetworkId:  "",
			UserId:     [8]byte{},
			RemoteAddr: nil,
			SrcIP:      nil,
			DstIP:      nil,
			FrameType:  0,
			Peer:       nil,
			Encrypt:    true,
		}
		return frame
	})

	buffPool = NewPool(func() any {
		buff := make([]byte, chacha20poly1305.NonceSize)
		return &buff
	})
	return
}

func NewPool(new func() any) *MemoryPool {
	return &MemoryPool{
		pool: sync.Pool{New: new},
	}
}

func (p *MemoryPool) Get() any {
	//p.lock.Lock()
	return p.pool.Get()
}

func (p *MemoryPool) Put(x any) {
	//defer p.lock.Unlock()
	p.pool.Put(x)
}

func (n *Node) GetBuffer() *[chacha20poly1305.NonceSize]byte {
	return n.pools.buffPool.Get().(*[chacha20poly1305.NonceSize]byte)
}

func (n *Node) PutBuffer(buffPtr *[chacha20poly1305.NonceSize]byte) {
	n.pools.buffPool.Put(buffPtr)
}

func (n *Node) GetFrame() *Frame {
	frame := n.pools.framePool.Get().(*Frame)
	frame.ST = time.Now()
	return frame
}

func (n *Node) PutFrame(framePtr *Frame) {
	framePtr.Size = 0
	framePtr.NetworkId = ""
	framePtr.SrcIP = nil
	framePtr.DstIP = nil
	framePtr.Packet = framePtr.Buff
	framePtr.Peer = nil
	framePtr.RemoteAddr = nil
	framePtr.FrameType = 0
	n.pools.framePool.Put(framePtr)
}
