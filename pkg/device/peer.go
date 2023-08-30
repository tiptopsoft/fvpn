// Copyright 2023 Tiptopsoft, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//	http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package device

import (
	"github.com/tiptopsoft/fvpn/pkg/device/conn"
	"github.com/tiptopsoft/fvpn/pkg/packet"
	"github.com/tiptopsoft/fvpn/pkg/packet/handshake"
	"github.com/tiptopsoft/fvpn/pkg/security"
	"github.com/tiptopsoft/fvpn/pkg/tun"
	"github.com/tiptopsoft/fvpn/pkg/util"
	"net"
	"sync"
	"sync/atomic"
	"time"
)

type Peer struct {
	isTry       atomic.Bool
	node        *Node
	mode        int
	ip          string
	device      tun.Device
	isRelay     bool
	index       atomic.Int32
	st          time.Time
	keepaliveCh chan int //1：exit keepalive 2: exit send packet 3 exit timer
	sendCh      chan int
	checkCh     chan int
	p2p         bool
	lock        sync.Mutex
	status      bool
	pubKey      security.NoisePublicKey
	bind        conn.Interface
	endpoint    conn.Endpoint //
	cache       Interface

	queue struct {
		outBound *OutBoundQueue // data to write to dst Peer
		inBound  *InBoundQueue  // data write to tun
	}
	cipher security.Codec
}

func (p *Peer) GetIP() string {
	return p.ip
}

func (p *Peer) GetCodec() security.Codec {
	return p.cipher
}

func (p *Peer) SetCodec(cipherFunc security.Codec) {
	p.cipher = cipherFunc
}

func (p *Peer) GetStatus() bool {
	return p.status
}

func (p *Peer) SetStatus(status bool) {
	p.status = status
}

func (p *Peer) GetEndpoint() conn.Endpoint {
	return p.endpoint
}

func (p *Peer) SetP2P(p2p bool) {
	p.p2p = p2p
}

func (p *Peer) GetP2P() bool {
	return p.p2p
}

func (p *Peer) SetMode(mode int) {
	p.mode = mode
}

func (p *Peer) Start() {
	p.lock.Lock()
	defer p.lock.Unlock()
	p.SetStatus(true)
	if p.mode == 1 {
		go func() {
			timer := time.NewTimer(time.Second * 0)
			defer timer.Stop()
			for {
				select {
				case <-p.keepaliveCh:
					return
				case <-timer.C:
					p.handshake(net.ParseIP(p.ip))
					p.keepalive()
					timer.Reset(time.Second * 10)
				}
			}
		}()

		go func() {
			timer := time.NewTimer(time.Second * 30)
			defer timer.Stop()
			for {
				select {
				case <-p.checkCh:
					return
				case <-timer.C:
					b := p.check()
					if b {
						//shutdown this Peer
						p.close()
						logger.Warnf("build p2p to node [%v] failed,exit now", p.ip)
					}
					timer.Reset(time.Second * 30)
				}
			}
		}()
	}
}

func (p *Peer) SetEndpoint(ep conn.Endpoint) {
	p.endpoint = ep
}

func (p *Peer) handshake(dstIP net.IP) {
	hpkt := handshake.NewPacket(util.HandShakeMsgType, util.UCTL.UserId)
	hpkt.Header.SrcIP = p.node.device.Addr()
	hpkt.Header.DstIP = dstIP
	hpkt.PubKey = p.pubKey
	buff, err := handshake.Encode(hpkt)
	if err != nil {
		return
	}

	size := len(buff)
	f := NewFrame()
	copy(f.Packet[:size], buff)
	f.Size = size
	f.FrameType = util.HandShakeMsgType
	f.Peer = p
	logger.Debugf("sending handshake pubkey to: %v, pubKey: %v, remote address: [%v], type: [%v]", dstIP.String(), p.pubKey, p.GetEndpoint().DstToString(), util.GetFrameTypeName(util.HandShakeMsgType))
	p.node.PutPktToOutbound(f)
}

func (p *Peer) PutPktToOutbound(pkt *Frame) {
	p.queue.outBound.c <- pkt
}

func (p *Peer) keepalive() {
	pkt, err := packet.NewHeader(util.KeepaliveMsgType, "")
	if err != nil {
		return
	}
	buff, err := packet.Encode(pkt)
	if err != nil {
		return
	}
	size := len(buff)
	f := NewFrame()
	f.Peer = p
	copy(f.Packet[:size], buff)
	f.Size = size
	f.FrameType = util.KeepaliveMsgType

	p.node.PutPktToOutbound(f)
}

func (p *Peer) check() bool {
	if p.isRelay || p.p2p {
		return false
	}
	st := time.Since(p.st)
	if st.Seconds() >= 30 {
		return true
	}

	return false
}

func (p *Peer) close() {
	p.checkCh <- 1
	p.sendCh <- 1
	p.keepaliveCh <- 1
	p.status = false
	p.isTry.Store(false)
	p.cache.Set(util.UCTL.UserId, p.ip, p)
	logger.Debugf("================Peer stop signal have send to Peer: %v", p.GetEndpoint().DstToString())
}
