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
	"github.com/tiptopsoft/fvpn/pkg/nets"
	"github.com/tiptopsoft/fvpn/pkg/packet"
	"github.com/tiptopsoft/fvpn/pkg/packet/handshake"
	"github.com/tiptopsoft/fvpn/pkg/packet/peer"
	"github.com/tiptopsoft/fvpn/pkg/packet/register/ack"
	"github.com/tiptopsoft/fvpn/pkg/security"
	"github.com/tiptopsoft/fvpn/pkg/util"
)

type Handler interface {
	Handle(ctx context.Context, frame *Frame) error
}

type HandlerFunc func(context.Context, *Frame) error

func (f HandlerFunc) Handle(ctx context.Context, frame *Frame) error {
	return f(ctx, frame)
}

type Middleware func(Handler) Handler

// Chain wrap middleware in order execute
func Chain(middlewares ...Middleware) func(Handler) Handler {
	return func(h Handler) Handler {
		for i := len(middlewares) - 1; i >= 0; i-- {
			h = middlewares[i](h)
		}

		return h
	}
}

func WithMiddlewares(handler Handler, middlewares ...Middleware) Handler {
	return Chain(middlewares...)(handler)
}

func (n *Node) tunInHandler() HandlerFunc {
	return func(ctx context.Context, frame *Frame) error {
		//defer frame.Unlock()
		n.PutPktToOutbound(frame)
		return nil
	}
}

// Handle union udp handler
func (n *Node) udpInHandler() HandlerFunc {
	return func(ctx context.Context, frame *Frame) error {
		//dest := ctx.Value("destAddr").(string)
		buff := frame.Packet[:]
		headerBuff, err := packet.Decode(buff)
		if err != nil {
			return err
		}

		//frame.FrameType = headerBuff.Flags
		switch headerBuff.Flags {
		case util.MsgTypeRegisterAck:
			regAck, err := ack.Decode(buff)
			if err != nil {
				return err
			}
			logger.Debugf("register success, got server server ack: (%v)", regAck)
		case util.MsgTypeQueryPeer:
			n.handleQueryPeers(frame)
		case util.MsgTypePacket:
			n.PutPktToInbound(frame)
		case util.HandShakeMsgType:
			//cache dst peer when receive a handshake
			headerPkt, err := handshake.Decode(util.HandShakeMsgType, frame.Buff)
			if err != nil {
				logger.Errorf("invalid handshake packet: %v", err)
				return err
			}

			p := NewPeer(util.UCTL.UserId, frame.SrcIP.String(), headerPkt.PubKey, n.cache, n)
			ep := nets.NewEndpoint(frame.RemoteAddr.String())
			p.SetEndpoint(ep)
			n.cache.SetPeer(frame.UidString(), frame.SrcIP.String(), p)
			//}
			//build handshake resp
			hpkt := handshake.NewPacket(util.HandShakeMsgTypeAck, frame.UidString())
			hpkt.Header.SrcIP = frame.DstIP
			hpkt.Header.DstIP = frame.SrcIP
			hpkt.PubKey = n.privateKey.NewPubicKey()
			buff, err := handshake.Encode(hpkt)
			if err != nil {
				return err
			}

			newFrame := NewFrame()
			newFrame.Size = len(buff)
			newFrame.Peer = p
			newFrame.UserId = frame.UserId
			newFrame.FrameType = util.HandShakeMsgTypeAck
			newFrame.DstIP = frame.SrcIP
			copy(newFrame.Packet[:newFrame.Size], buff)
			n.PutPktToOutbound(newFrame)
		case util.HandShakeMsgTypeAck: //use for relay
			//err = n.cache.SetPeer(frame.UidString(), frame.SrcIP.String(), p)
			_, err = CachePeers(n.privateKey, frame, n.cache, n)
			p, err := n.cache.GetPeer(frame.UidString(), frame.SrcIP.String())
			p.p2p = true
			n.cache.SetPeer(frame.UidString(), frame.SrcIP.String(), p)
			if err != nil {
				return err
			}
		case util.KeepaliveMsgType:
		}

		return nil
	}
}

func CachePeers(privateKey security.NoisePrivateKey, frame *Frame, cache Interface, node *Node) (*Peer, error) {
	hpkt, err := handshake.Decode(util.HandShakeMsgTypeAck, frame.Buff)
	if err != nil {
		logger.Errorf("invalid handshake packet: %v", err)
		return nil, err
	}
	uid := frame.UidString()
	srcIP := frame.SrcIP.String()
	logger.Debugf("got remote peer: %v, pubKey: %v", srcIP, hpkt.PubKey)

	p := NewPeer(uid, srcIP, hpkt.PubKey, cache, node)
	p.node = node
	ep := nets.NewEndpoint(frame.RemoteAddr.String())
	p.SetEndpoint(ep)
	p.cipher = security.NewCipher(privateKey, hpkt.PubKey)
	err = cache.SetPeer(uid, srcIP, p)
	p.start()

	if err != nil {
		return nil, err
	}

	return p, nil
}

func (n *Node) handleQueryPeers(frame *Frame) {
	peers, _ := peer.Decode(frame.Packet[:])
	logger.Debugf("go peers from remote: %v", peers.Peers)
	for _, info := range peers.Peers {
		ip := info.IP
		if ip.String() == n.device.IPToString() {
			continue
		}

		addr := info.RemoteAddr
		p := NewPeer(frame.UidString(), ip.String(), security.NoisePublicKey{}, n.cache, n) //now has no pubKey
		if p.endpoint == nil {
			p.SetEndpoint(nets.NewEndpoint(addr.String()))
		} else if p.endpoint.DstToString() != addr.String() {
			p.SetEndpoint(nets.NewEndpoint(addr.String()))
		}
		err := n.cache.SetPeer(frame.UidString(), ip.String(), p)
		if err != nil {
			return
		}

		p.start()
		if p.status {
			p.handshake(ip)
		}
	}
}
