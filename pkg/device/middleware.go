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
	"fmt"
	"github.com/tiptopsoft/fvpn/pkg/packet"
	"github.com/tiptopsoft/fvpn/pkg/util"
)

// AuthCheck Middleware auth handler to check user login, if not, return an error tell user to login first.
func AuthCheck() func(handler Handler) Handler {
	return func(next Handler) Handler {
		return HandlerFunc(func(ctx context.Context, frame *Frame) error {
			//username, password, err := util.GetUserInfo()
			//if err != nil {
			//	return err
			//}
			//
			//client := http.New("http://211.159.225.186:443")
			//req := new(http.LoginRequest)
			//req.username = username
			//req.password = password
			//loginResp, err := client.Login(*req)
			//if err != nil {
			//	return errors.New("user should login first")
			//}
			//
			//if loginResp.Token == "" {
			//	return errors.New("token is nil, please login again")
			//}

			return next.Handle(ctx, frame)
		})
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
				decoded, err := peer.GetCodec().Decode(buff)
				if err != nil {
					return err
				}
				frame.Size = len(decoded) + offset
				copy(frame.Packet[offset:frame.Size], decoded)
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
				encoded, err := peer.GetCodec().Encode(buff)
				if err != nil {
					return err
				}
				frame.Size = offset + len(encoded)
				copy(frame.Packet[offset:frame.Size], encoded)
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
				b := n.netCtl.Access(frame.UserIdString(), ip)
				if !b {
					return fmt.Errorf("no access to destination: [%v]", ip)
				}
			}
			return next.Handle(ctx, frame)
		})
	}
}
