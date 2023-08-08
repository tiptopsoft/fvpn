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
	"github.com/tiptopsoft/fvpn/pkg/nets"
	"github.com/tiptopsoft/fvpn/pkg/packet"
	"github.com/tiptopsoft/fvpn/pkg/packet/handshake"
	"github.com/tiptopsoft/fvpn/pkg/security"
	"github.com/tiptopsoft/fvpn/pkg/util"
	"go.uber.org/atomic"
	"net"
	"sync"
	"time"
)

// Peer a destination will have a peer in fvpn, can connect to each other.
// a RegServer also is a peer
type Peer struct {
	id          uint64
	isRelay     bool
	index       atomic.Int32
	st          time.Time
	keepaliveCh chan int //1：exit keepalive 2: exit send packet 3 exit timer
	sendCh      chan int
	checkCh     chan int
	p2p         bool
	lock        sync.Mutex
	status      bool
	PubKey      security.NoisePublicKey
	node        *Node
	endpoint    nets.Endpoint //

	queue struct {
		outBound *OutBoundQueue // data to write to dst peer
		inBound  *InBoundQueue  // data write to tun
	}
	cipher security.CipherFunc
}

func (p *Peer) GetCodec() security.CipherFunc {
	return p.cipher
}

func (p *Peer) start() {
	//p.lock.Lock()
	//defer p.lock.Unlock()
	if p.index.Load() > 3 {
		logger.Debugf("peer %v have try too much times", p.endpoint.DstToString())
		return
	}
	p.index.Inc()
	if p.status == true {
		logger.Debugf("peer has started: %v", p.endpoint.DstToString())
		return
	} else {
		logger.Debugf("peer starting......")
	}

	p.status = true
	if p.node != nil && p.node.mode == 1 {
		p.PubKey = p.node.privateKey.NewPubicKey() //use to send to remote for exchange pubKey
		go p.SendPackets()
		//go p.WriteToDevice()

		go func() {
			timer := time.NewTimer(time.Second * 10)
			defer timer.Stop()
			for {
				select {
				case <-p.keepaliveCh:
					return
				case <-timer.C:
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
						//shutdown this peer
						p.close()
					}
					timer.Reset(time.Second * 30)
				}
			}
		}()
	}
}

func (p *Peer) SetEndpoint(ep nets.Endpoint) {
	p.endpoint = ep
}

func (p *Peer) GetEndpoint() nets.Endpoint {
	return p.endpoint
}

func NewPeer(uid, srcIP string, pk security.NoisePublicKey, cache Interface, node *Node) *Peer {
	logger.Debugf("will create peer for userId: %v, ip: %v", uid, srcIP)
	peer, _ := cache.GetPeer(uid, srcIP)
	if peer != nil {
		return peer
	}

	p := new(Peer)
	p.id = uint64(time.Now().Nanosecond())
	p.st = time.Now()
	p.checkCh = make(chan int, 1)
	p.sendCh = make(chan int, 1)
	p.keepaliveCh = make(chan int, 1)
	p.PubKey = pk
	p.queue.outBound = NewOutBoundQueue()
	p.queue.inBound = NewInBoundQueue()
	if node != nil {
		p.node = node
	}

	cache.SetPeer(uid, srcIP, p)
	logger.Debugf("created peer for : %v", srcIP)
	return p
}

func (p *Peer) handshake(dstIP net.IP) {
	hpkt := handshake.NewPacket(util.HandShakeMsgType, util.UCTL.UserId)
	hpkt.Header.SrcIP = p.node.device.Addr()
	hpkt.Header.DstIP = dstIP
	hpkt.PubKey = p.PubKey
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
	logger.Debugf("sending handshake pubkey to: %v, pubKey: %v, remote address: [%v], type: [%v]", dstIP.String(), p.PubKey, p.endpoint.DstToString(), util.GetFrameTypeName(util.HandShakeMsgType))
	p.PutPktToOutbound(f)
}

func (p *Peer) PutPktToOutbound(pkt *Frame) {
	p.queue.outBound.c <- pkt
	logger.Debugf("Adding pkt to peer: [%v], data type: [%v]", p.id, util.GetFrameTypeName(pkt.FrameType))
}

func (p *Peer) SendPackets() {
	for {
		select {
		case <-p.sendCh:
			return
		case pkt := <-p.queue.outBound.c:
			send, err := p.node.net.bind.Send(pkt.Packet[:pkt.Size], p.endpoint)
			if err != nil {
				logger.Error(err)
				continue
			}
			logger.Debugf("node [%v] has send %d packets to %s, data type: [%v]", p.id, send, p.endpoint.DstToString(), util.GetFrameTypeName(pkt.FrameType))
		default:

		}
	}
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

	p.PutPktToOutbound(f)
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
	logger.Debug("================peer stop signal have send to peer: %v", p.endpoint.DstToString())
}
