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
			logger.Debugf("got handshake msg type, data: %v", frame.Packet[:frame.Size])
			headerPkt, err := handshake.Decode(frame.Buff)
			if err != nil {
				logger.Errorf("invalid handshake packet: %v", err)
				return err
			}
			p, err := n.cache.GetPeer(util.UCTL.UserId, frame.SrcIP.String())
			if err != nil || p == nil {
				p = n.NewPeer(headerPkt.PubKey)
				p.node = n
				ep := nets.NewEndpoint(frame.RemoteAddr.String())
				p.SetEndpoint(ep)
				n.cache.SetPeer(frame.UidString(), frame.SrcIP.String(), p)
			}
			//build handshake resp
			hPktack := handshake.NewPacket(util.HandShakeMsgTypeAck, frame.UidString())
			logger.Debugf("got packet srcIP: %v, dstIP: %v, data type: [%v]", frame.SrcIP, frame.DstIP, util.GetFrameTypeName(util.HandShakeMsgType))
			hPktack.Header.SrcIP = frame.DstIP //dstIP = 2
			hPktack.Header.DstIP = frame.SrcIP //srcIP = 1
			hPktack.PubKey = n.privateKey.NewPubicKey()
			buff, err := handshake.Encode(hPktack)
			if err != nil {
				return err
			}

			frame.Packet = buff
			frame.Size = len(buff)
			frame.DstIP = frame.SrcIP //dstIP = 1
			frame.Peer = p
			logger.Debugf(">>>>>>>will send a handshakd ack to remote back, dst: [%v]", frame.Peer.endpoint.DstToString())
			n.PutPktToOutbound(frame)
		case util.HandShakeMsgTypeAck: //use for relay
			//cache dst peer when receive a handshake
			logger.Debugf("got handshake msg type in handshake ack, data: %v, data type: [%v]", frame.Packet[:frame.Size], util.GetFrameTypeName(util.HandShakeMsgTypeAck))
			//err = n.cache.SetPeer(frame.UidString(), frame.SrcIP.String(), p)
			_, err = CachePeerToLocal(n.privateKey, frame, n.cache, n)
			p, err := n.cache.GetPeer(frame.UidString(), frame.SrcIP.String())
			p.p2p = true
			n.cache.SetPeer(frame.UidString(), frame.SrcIP.String(), p)
			if err != nil {
				return err
			}
		case util.KeepaliveMsgType:
			logger.Debugf("got keepalived packets from :%v, data: %v", frame.RemoteAddr, frame.Packet[:frame.Size])
		}

		return nil
	}
}

func CachePeerToLocal(privateKey security.NoisePrivateKey, frame *Frame, cache CacheFunc, node *Node) (*Peer, error) {
	hpkt, err := handshake.Decode(frame.Buff)
	if err != nil {
		logger.Errorf("invalid handshake packet: %v", err)
		return nil, err
	}
	uid := frame.UidString()
	srcIP := frame.SrcIP.String()
	logger.Debugf("got remote peer: %v, pubKey: %v", srcIP, hpkt.PubKey)

	p, err := cache.GetPeer(uid, srcIP)
	if err != nil || p == nil {
		p = node.NewPeer(hpkt.PubKey)
	}
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
		p, err := n.cache.GetPeer(frame.UidString(), ip.String())
		if err != nil || p == nil {
			p = n.NewPeer(security.NoisePublicKey{}) //now has no pubKey
			p.SetEndpoint(nets.NewEndpoint(addr.String()))
			err = n.cache.SetPeer(frame.UidString(), ip.String(), p)
		} else {
			if p.endpoint.DstToString() != addr.String() {
				p.SetEndpoint(nets.NewEndpoint(addr.String()))
			}
		}

		if err != nil {
			return
		}
		p.start()
		if p.status {
			p.handshake(ip)
		}
	}
}
