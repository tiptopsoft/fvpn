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
	"fmt"
	"github.com/tiptopsoft/fvpn/pkg/packet"
	"github.com/tiptopsoft/fvpn/pkg/util"
)

func Decode() func(Handler) Handler {
	return func(next Handler) Handler {
		return HandlerFunc(func(ctx context.Context, frame *Frame) error {
			if frame.FrameType == util.MsgTypePacket {
				offset := packet.HeaderBuffSize
				buff := frame.Packet[offset:frame.Size]
				//cache := ctx.Value("cache").(CacheFunc)
				peer := frame.GetPeer()
				if peer == nil {
					return fmt.Errorf("dst ip: %v peer not found", frame.DstIP.String())
				}

				logger.Debugf("use src peer: [%v] to decode", peer.endpoint.DstIP().String())

				logger.Debugf("data before decode: %v", buff)
				decoded, err := peer.GetCodec().Decode(buff)
				if err != nil {
					return err
				}
				frame.Clear()
				copy(frame.Packet[0:offset], frame.Buff[0:offset])
				copy(frame.Packet[offset:], decoded)
				frame.Size = len(decoded) + offset
				logger.Debugf("data after decode: %v", frame.Packet[:frame.Size])
			}
			return next.Handle(ctx, frame)
		})
	}
}

// Encode Middleware encrypt use exchangeKey
func Encode() func(Handler) Handler {
	return func(next Handler) Handler {
		return HandlerFunc(func(ctx context.Context, frame *Frame) error {
			if frame.FrameType == util.MsgTypePacket {
				offset := packet.HeaderBuffSize
				buff := frame.Packet[offset:frame.Size]
				//cache := ctx.Value("cache").(CacheFunc)

				//peer, err := cache.GetPeer(UCTL.UserId, frame.DstIP.String())
				//if err != nil  {
				//	return errors.New("peer not found, if you want to use relay, please to put relay true")
				//}

				peer := frame.GetPeer()
				logger.Debugf("peer is :%v, data before encode: %v", peer.GetEndpoint().DstIP(), buff)
				encoded, err := peer.GetCodec().Encode(buff)
				if err != nil {
					return err
				}
				copy(frame.Packet[offset:], encoded)
				frame.Size = offset + len(encoded)
				logger.Debugf("data after encode: %v", frame.Packet[:frame.Size])
				if err != nil {
					return err
				}
			}
			return next.Handle(ctx, frame)
		})
	}
}

// AllowNetwork valid user can join a network or a node, so here will check
func (n *Node) AllowNetwork() func(Handler) Handler {
	return func(next Handler) Handler {
		return HandlerFunc(func(ctx context.Context, frame *Frame) error {
			if frame.FrameType == util.MsgTypePacket {
				ip := frame.DstIP.String()
				b := n.netCtl.Access(frame.UidString(), ip)
				if !b {
					return fmt.Errorf("has no access to IP: %v", ip)
				}
			}
			return next.Handle(ctx, frame)
		})
	}
}
