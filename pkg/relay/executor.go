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

package relay

import (
	"context"
	"fmt"
	"github.com/tiptopsoft/fvpn/pkg/device"
	"github.com/tiptopsoft/fvpn/pkg/device/conn"
	"github.com/tiptopsoft/fvpn/pkg/packet/handshake"
	"github.com/tiptopsoft/fvpn/pkg/packet/peer"
	"github.com/tiptopsoft/fvpn/pkg/security"
	"github.com/tiptopsoft/fvpn/pkg/util"
	"net"
)

func (r *RegServer) ReadFromUdp() {
	logger.Infof("start a udp loop")
	for {
		buffPtr := r.pool.Get()

		ctx := context.Background()
		ctx = context.WithValue(ctx, "cache", r.cache)
		frame := device.NewFrame()
		frame.Packet = *buffPtr
		frame.Ctx = ctx
		n, addr, err := r.conn.ReadFromUDP(frame.Packet[:])
		if err != nil || n < 0 {
			logger.Error("no data exists")
			continue
		}
		packetHeader, err := util.GetPacketHeader(frame.Packet[:])
		if err != nil {
			logger.Errorf("get header falied. %v", err)
			continue
		}

		logger.Debugf("Read from %v udp %d byte, srcIP: %v, dstIP: %v, data type: [%v]", addr, n, packetHeader.SrcIP, packetHeader.DstIP, util.GetFrameTypeName(packetHeader.Flags))
		frame.Size = n
		frame.FrameType = packetHeader.Flags
		frame.RemoteAddr = addr
		frame.SrcIP = packetHeader.SrcIP
		frame.DstIP = packetHeader.DstIP
		frame.UserId = packetHeader.UserId
		//decode use origin peer
		frame.Peer, _ = r.cache.Get(frame.UidString(), frame.SrcIP.String())
		r.PutPktToInbound(frame)
	}
}

func (r *RegServer) writeUdpHandler() device.HandlerFunc {
	return func(ctx context.Context, pkt *device.Frame) error {
		n, err := r.conn.WriteToUDP(pkt.Packet[:pkt.Size], pkt.RemoteAddr)
		if err != nil {
			return err
		}
		logger.Debugf("registry write %d size to %v, data: %v", n, pkt.RemoteAddr, pkt.Packet[:pkt.Size])
		return nil
	}
}

// serverUdpHandler  core self handler
func (r *RegServer) serverUdpHandler() device.HandlerFunc {
	return func(ctx context.Context, frame *device.Frame) error {
		switch frame.FrameType {
		//case util.MsgTypeRegisterSuper:
		//	err := r.register(frame)
		//	h, err := packet.NewHeader(util.MsgTypeRegisterAck, frame.NetworkId)
		//	if err != nil {
		//		logger.Errorf("build resp failed. err: %v", err)
		//	}
		//	f, _ := packet.Encode(h)
		//	frame.Packet = f
		//	break
		case util.MsgTypePacket:
			p, err := r.cache.Get(frame.UidString(), frame.DstIP.String())
			if err != nil || p == nil {
				return fmt.Errorf("peer %v is not found", frame.DstIP.String())
			}
			logger.Debugf("write packet to peer %v: ", p)
			frame.RemoteAddr = p.GetEndpoint().DstIP()
			frame.Peer = p //change peer to dst peer
			r.PutPktToOutbound(frame)
		case util.MsgTypeQueryPeer:
			peers := r.cache.Lists(frame.UidString())
			peerAck := peer.NewPacket(frame.UidString())

			if len(peers) > 0 {
				for ip, value := range peers {
					info := peer.Info{
						IP:         net.ParseIP(ip),
						RemoteAddr: *value.Peer.GetEndpoint().DstIP(),
					}
					peerAck.Peers = append(peerAck.Peers, info)
				}
			}

			buff, _ := peer.Encode(peerAck)
			r.pool.Put(&frame.Packet)

			newFrame := device.NewFrame()
			buffPtr := r.pool.Get()
			newFrame.Packet = *buffPtr
			copy(newFrame.Packet[:], buff)
			newFrame.UserId = frame.UserId
			newFrame.RemoteAddr = frame.RemoteAddr
			newFrame.FrameType = util.MsgTypeQueryPeer
			newFrame.Size = len(buff)
			r.PutPktToOutbound(newFrame)
			r.pool.Put(buffPtr)
		case util.HandShakeMsgType:
			srcIP := frame.SrcIP.String()
			uid := frame.UidString()
			pktBuf, err := handshake.Decode(util.HandShakeMsgTypeAck, frame.Packet[:])
			if err != nil {
				return err
			}
			node := new(device.Node)
			p := node.NewPeer(uid, srcIP, r.key.privateKey.NewPubicKey(), r.cache)
			p.SetEndpoint(conn.NewEndpoint(frame.RemoteAddr.String()))
			p.SetCodec(security.New(r.key.privateKey, pktBuf.PubKey))
			logger.Infof("registry build p2p with node: %v success", srcIP)
			r.cache.Set(uid, srcIP, p)
			//build handshakeAck resp
			pkt := handshake.NewPacket(util.HandShakeMsgTypeAck, uid)
			pkt.Header.SrcIP = frame.DstIP
			pkt.Header.DstIP = frame.SrcIP
			pkt.PubKey = r.key.privateKey.NewPubicKey()
			buff, err := handshake.Encode(pkt)
			if err != nil {
				return err
			}

			newFrame := device.NewFrame()
			newFrame.Size = len(buff)
			newFrame.RemoteAddr = frame.RemoteAddr
			copy(newFrame.Packet[:newFrame.Size], buff)
			r.PutPktToOutbound(newFrame)
		}

		return nil
	}
}

func (r *RegServer) register(frame *device.Frame) (err error) {
	p := new(device.Peer)
	ep := conn.NewEndpoint(frame.RemoteAddr.String())
	p.SetEndpoint(ep)
	err = r.cache.Set(frame.UidString(), frame.SrcIP.String(), p)
	return
}
