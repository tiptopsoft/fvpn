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
	"github.com/tiptopsoft/fvpn/pkg/packet"
	"github.com/tiptopsoft/fvpn/pkg/packet/handshake"
	"github.com/tiptopsoft/fvpn/pkg/packet/peer"
	"github.com/tiptopsoft/fvpn/pkg/packet/register/ack"
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

			p := NewPeer(util.UCTL.UserId, frame.SrcIP.String(), headerPkt.PubKey, n.cache, n.mode, n.net.bind, n.device)

			if p.GetEndpoint() == nil || p.GetEndpoint().DstToString() != frame.RemoteAddr.String() {
				//更新peer
				ep := NewEndpoint(frame.RemoteAddr.String())
				p.SetEndpoint(ep)
				n.cache.SetPeer(frame.UidString(), frame.SrcIP.String(), p)
			}

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
			_, err = CachePeers(n.privateKey, frame, n.cache, n.mode, n.net.bind, n.device)
			p, err := n.cache.GetPeer(frame.UidString(), frame.SrcIP.String())
			p.SetP2P(true)
			n.cache.SetPeer(frame.UidString(), frame.SrcIP.String(), p)
			if err != nil {
				return err
			}
		case util.KeepaliveMsgType:
		}

		return nil
	}
}

func (n *Node) handleQueryPeers(frame *Frame) {
	peers, _ := peer.Decode(frame.Packet[:])
	logger.Debugf("got peers from remote: %v", peers.Peers)
	for _, info := range peers.Peers {
		ip := info.IP
		if ip.String() == n.device.IPToString() {
			continue
		}

		addr := info.RemoteAddr
		p := NewPeer(frame.UidString(), ip.String(), n.privateKey.NewPubicKey(), n.cache, n.mode, n.net.bind, n.device) //now has no pubKey
		if p.GetEndpoint() == nil {
			p.SetEndpoint(NewEndpoint(addr.String()))
		}
		logger.Debugf("print now endpoint: [%v]", p.GetEndpoint().DstToString())
		err := n.cache.SetPeer(frame.UidString(), ip.String(), p)
		if err != nil {
			return
		}

		p.Start()
		if p.GetStatus() {
			p.Handshake(ip)
		}
	}
}
