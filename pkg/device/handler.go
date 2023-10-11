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
	"context"
	"github.com/tiptopsoft/fvpn/pkg/device/conn"
	"github.com/tiptopsoft/fvpn/pkg/packet"
	"github.com/tiptopsoft/fvpn/pkg/packet/handshake"
	"github.com/tiptopsoft/fvpn/pkg/packet/peer"
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
		n.PutPktToOutbound(frame)
		return nil
	}
}

// Handle union udp handler
func (n *Node) udpInHandler() HandlerFunc {
	return func(ctx context.Context, frame *Frame) error {
		buff := frame.Packet[:]
		headerBuff, err := packet.Decode(buff)
		if err != nil {
			return err
		}

		switch headerBuff.Flags {
		case util.MsgTypeQueryPeer:
			n.handleQueryPeers(frame)
		case util.MsgTypePacket:
			n.PutPktToInbound(frame)
		case util.HandShakeMsgType:
			//cache dst Peer when receive a handshake
			headerPkt, err := handshake.Decode(util.HandShakeMsgType, frame.Packet[:])
			if err != nil {
				logger.Errorf("invalid handshake packet: %v", err)
				return err
			}

			p := n.NewPeer(util.Info().GetUserId(), frame.SrcIP.String(), headerPkt.PubKey, n.cache)
			p.node = n

			//if just one node behind Symmetric nat, also update endpoint to build p2p
			if p.GetEndpoint() == nil || p.GetEndpoint().DstToString() != frame.RemoteAddr.String() {
				logger.Debugf("this is a symetric nat device: %s", frame.RemoteAddr.String())
				//更新peer
				ep := conn.NewEndpoint(frame.RemoteAddr.String())
				p.SetEndpoint(ep)
				n.cache.Set(frame.UidString(), frame.SrcIP.String(), p)
			}

			//build handshake ack
			pkt := handshake.NewPacket(util.HandShakeMsgTypeAck, frame.UidString())
			pkt.Header.SrcIP = frame.DstIP
			pkt.Header.DstIP = frame.SrcIP
			pkt.PubKey = n.privateKey.NewPubicKey()

			newFrame := n.GetFrame()
			if newFrame.Size, err = pkt.Encode(newFrame.Packet[:]); err != nil {
				return err
			}

			newFrame.Peer = p
			newFrame.UserId = frame.UserId
			newFrame.FrameType = util.HandShakeMsgTypeAck
			newFrame.DstIP = frame.SrcIP
			n.PutFrame(frame)
			p.sendBuffer(newFrame, newFrame.GetPeer().GetEndpoint())
			n.PutFrame(newFrame)
		case util.HandShakeMsgTypeAck:
			srcIP := frame.SrcIP.String()
			uid := frame.UidString()
			p, err := n.cache.Get(uid, srcIP)
			pkt, err := handshake.Decode(util.HandShakeMsgTypeAck, frame.Packet[:])
			if err != nil {
				return err
			}
			ep := conn.NewEndpoint(frame.RemoteAddr.String())
			p.SetEndpoint(ep)
			p.SetCodec(security.New(n.privateKey, pkt.PubKey))
			if !p.p2p {
				p.SetP2P(true)
				logger.Infof("node [%v] build a p2p to node [%v]", frame.DstIP, frame.SrcIP)
			}
			err = n.cache.Set(uid, srcIP, p)
			if !p.GetStatus() {
				p.Start()
			}
			n.cache.Set(uid, srcIP, p)
			if err != nil {
				return err
			}
			n.PutFrame(frame)
		}

		return nil
	}
}

func (n *Node) handleQueryPeers(frame *Frame) {
	defer n.PutFrame(frame)
	peers, _ := peer.Decode(frame.Packet[:])
	logger.Debugf("list peers from registry: %v", peers.Peers)
	for _, info := range peers.Peers {
		dstIP := info.IP.String()
		uid := frame.UidString()
		if dstIP == n.device.IPToString() {
			//go over if dstIP is local dstIP
			continue
		}

		addr := info.RemoteAddr.String()
		p := n.NewPeer(uid, dstIP, n.privateKey.NewPubicKey(), n.cache)
		p.SetEndpoint(conn.NewEndpoint(addr))

		p.mode = 1
		if !p.GetStatus() {
			p.Start()
		}
	}
}
