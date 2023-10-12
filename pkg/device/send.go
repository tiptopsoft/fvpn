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
	"errors"
	"github.com/tiptopsoft/fvpn/pkg/packet"
	"github.com/tiptopsoft/fvpn/pkg/util"
	"time"
)

/** OutBound flow
 *
 * flow of send
 * 1. read from tun(sync)
 * 2. encrypt(sync)
 * 3. sendto peer(sync)
 */

func (n *Node) ReadFromTun() {
	userId := util.Info().GetUserId()
	newHeader, _ := packet.NewHeader(util.MsgTypePacket, userId)
	for {
		ctx := context.Background()
		frame := n.GetFrame()
		frame.UserId = n.userId
		frame.FrameType = util.MsgTypePacket
		size, err := n.device.Read(frame.Packet[:], packet.HeaderBuffSize)
		frame.Size = size + packet.HeaderBuffSize
		if err != nil {
			logger.Error(err)
			n.PutFrame(frame)
			continue
		}
		ipHeader, err := util.GetIPFrameHeader(frame.Packet[packet.HeaderBuffSize:])
		if err != nil {
			logger.Error(err)
			n.PutFrame(frame)
			continue
		}
		dstIP := ipHeader.DstIP.String()
		if dstIP == n.device.Addr().String() {
			n.PutPktToInbound(frame)
			continue
		}
		logger.Debugf("node %s receive %d byte, srcIP: %v, dstIP: %v", n.device.Name(), size, ipHeader.SrcIP, ipHeader.DstIP)
		if n.cfg.Relay.Force {
			frame.Peer = n.relay
		} else {
			if peer, err := n.cache.Get(util.Info().GetUserId(), dstIP); err != nil || peer == nil {
				//drop peer is not online
				n.PutFrame(frame)
				continue
			} else if !peer.GetP2P() && n.cfg.EnableRelay() {
				frame.Peer = n.relay
			} else {
				frame.Peer = peer
			}
		}

		logger.Debugf("frame's Peer is :%v", frame.Peer.GetEndpoint().DstToString())
		frame.SrcIP = n.device.Addr()
		frame.DstIP = ipHeader.DstIP

		frame.UserId = newHeader.UserId
		newHeader.SrcIP = frame.SrcIP
		newHeader.DstIP = frame.DstIP
		if _, err = newHeader.Encode(frame.Packet[:]); err != nil {
			logger.Error(err)
			n.PutFrame(frame)
			continue
		}

		if !n.cfg.Encrypt.Enable {
			frame.Encrypt = false
		}

		err = n.tunHandler.Handle(ctx, frame)
		if err != nil {
			logger.Error(err)
			n.PutFrame(frame)
			continue
		}
	}
}

// Encode Middleware encrypt use exchangeKey
func Encode() func(Handler) Handler {
	return func(next Handler) Handler {
		return HandlerFunc(func(ctx context.Context, frame *Frame) error {
			if frame.FrameType == util.MsgTypePacket && frame.Encrypt {
				offset := packet.HeaderBuffSize
				buff := frame.Packet[offset:frame.Size]
				peer := frame.GetPeer()
				logger.Debugf("Peer is :%v, data before encode: %v", peer.GetEndpoint().DstIP(), buff)
				if peer.GetCodec() == nil {
					logger.Debugf("")
				}
				if peer.GetCodec() == nil {
					return errors.New("node has not built success yet")
				}
				if _, err := peer.GetCodec().Encode(buff); err != nil {
					return err
				}
				frame.Size = offset + len(buff)
				//copy(frame.Packet[offset:frame.Size], buff)
				logger.Debugf("data after encode: %v", frame.Packet[:frame.Size])
			}
			return next.Handle(ctx, frame)
		})
	}
}

func (n *Node) tunInHandler() HandlerFunc {
	return func(ctx context.Context, frame *Frame) error {
		n.PutPktToOutbound(frame)
		return nil
	}
}

func (n *Node) WriteToDevice() {
	for {
		select {
		case pkt := <-n.queue.inBound.c:
			if pkt.FrameType == util.MsgTypePacket {
				size, err := n.device.Write(pkt.Packet[:pkt.Size], packet.HeaderBuffSize)
				if err != nil {
					logger.Error(err)
					continue
				}

				t := time.Since(pkt.ST)
				logger.Debugf("node write %d byte to %s, cost: [%v]", size, n.device.Name(), t)
				n.PutFrame(pkt)
			}

		}
	}
}
