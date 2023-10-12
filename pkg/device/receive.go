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
	"fmt"
	"github.com/tiptopsoft/fvpn/pkg/device/conn"
	"github.com/tiptopsoft/fvpn/pkg/packet"
	"github.com/tiptopsoft/fvpn/pkg/packet/handshake"
	"github.com/tiptopsoft/fvpn/pkg/security"
	"github.com/tiptopsoft/fvpn/pkg/util"
	"time"
)

/**InBound flow
 *
 */

func (n *Node) ReadFromUdp() {
	logger.Debugf("start thread to handle udp packet")
	defer func() {
		logger.Debugf("udp thread exited")
	}()
	for {
		ctx := context.Background()
		frame := n.GetFrame()
		size, remoteAddr, err := n.net.conn.Conn().ReadFromUDP(frame.Packet[:])
		if err != nil {
			logger.Errorf("udp read remote failed, err: %v", err)
			n.PutFrame(frame)
			return
		}
		frame.Size = size
		frame.RemoteAddr = remoteAddr
		n.udpProcess(ctx, frame)
	}
}

func (n *Node) udpProcess(ctx context.Context, frame *Frame) {
	hpkt, err := util.GetPacketHeader(frame.Packet[:])
	if err != nil {
		logger.Error(err)
	}
	dataType := util.GetFrameTypeName(hpkt.Flags)
	if dataType == "" {
		//drop
		logger.Errorf("got unknown data. size: %d", frame.Size)
		n.PutFrame(frame)
		return
	}
	logger.Debugf("udp receive %d byte from %s, data type: [%v]", frame.Size, frame.RemoteAddr, dataType)

	frame.SrcIP = hpkt.SrcIP
	frame.DstIP = hpkt.DstIP
	frame.UserId = hpkt.UserId
	frame.FrameType = hpkt.Flags

	frame.Peer, err = n.cache.Get(frame.UserIdString(), frame.SrcIP.String())
	if err != nil || !frame.Peer.GetP2P() {
		frame.Peer = n.relay
	}

	if !n.cfg.Encrypt.Enable {
		frame.Encrypt = false
	}

	err = n.udpHandler.Handle(ctx, frame)
	if err != nil {
		logger.Errorf("udp handler error: %v", err)
		n.PutFrame(frame)
		return
	}
	dt := time.Since(frame.ST)
	logger.Debugf("udp receive process finished, dataType: [%v], cost: [%v]", dataType, dt)
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
				n.cache.Set(frame.UserIdString(), frame.SrcIP.String(), p)
			}

			//build handshake ack
			pkt := handshake.NewPacket(util.HandShakeMsgTypeAck, frame.UserIdString())
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
			uid := frame.UserIdString()
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

func Decode() func(Handler) Handler {
	return func(next Handler) Handler {
		return HandlerFunc(func(ctx context.Context, frame *Frame) error {
			if frame.FrameType == util.MsgTypePacket && frame.Encrypt {
				offset := packet.HeaderBuffSize
				buff := frame.Packet[offset:frame.Size]
				peer := frame.GetPeer()
				if peer == nil {
					return fmt.Errorf("dst ip: %v Peer not found", frame.DstIP.String())
				}

				logger.Debugf("use src Peer: [%v] to decode", peer.GetEndpoint().DstIP().String())

				logger.Debugf("data before decode: %v", buff)
				if _, err := peer.GetCodec().Decode(buff); err != nil {
					return err
				}
				frame.Size = len(buff) + offset
				//copy(frame.Packet[offset:frame.Size], decoded)
				logger.Debugf("data after decode: %v", frame.Packet[:frame.Size])
			}
			return next.Handle(ctx, frame)
		})
	}
}

func (n *Node) WriteToUDP() {
	for {
		select {
		case pkt := <-n.queue.outBound.c:
			peer := pkt.Peer
			send, err := n.net.conn.Send(pkt.Packet[:pkt.Size], pkt.Peer.GetEndpoint())
			if err != nil {
				logger.Error(err)
				continue
			}
			logger.Debugf("node has send [%v] packets to %s from p2p: [%v], data type: [%v]", send, peer.GetEndpoint().DstToString(), peer.p2p, util.GetFrameTypeName(pkt.FrameType))
			n.PutFrame(pkt)
		}
	}
}
