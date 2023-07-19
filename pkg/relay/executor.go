package relay

import (
	"context"
	"fmt"
	"github.com/topcloudz/fvpn/pkg/nets"
	"github.com/topcloudz/fvpn/pkg/node"
	"github.com/topcloudz/fvpn/pkg/packet"
	"github.com/topcloudz/fvpn/pkg/packet/handshake"
	"github.com/topcloudz/fvpn/pkg/packet/peer"
	"github.com/topcloudz/fvpn/pkg/util"
	"net"
)

func (r *RegServer) ReadFromUdp() {
	logger.Infof("start a udp loop")
	for {
		ctx := context.Background()
		ctx = context.WithValue(ctx, "cache", r.cache)
		frame := node.NewFrame()
		frame.Ctx = ctx
		n, addr, err := r.conn.ReadFromUDP(frame.Buff[:])
		if err != nil || n < 0 {
			logger.Error("no data exists")
			continue
		}
		copy(frame.Packet, frame.Buff)
		packetHeader, err := util.GetPacketHeader(frame.Buff[:])
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
		frame.Peer, _ = r.cache.GetPeer(frame.UidString(), frame.SrcIP.String())
		r.PutPktToInbound(frame)
	}
}

func (r *RegServer) writeUdpHandler() node.HandlerFunc {
	return func(ctx context.Context, pkt *node.Frame) error {
		n, err := r.conn.WriteToUDP(pkt.Packet[:pkt.Size], pkt.RemoteAddr)
		if err != nil {
			return err
		}
		logger.Debugf("registry write %d size to %v, data: %v", n, pkt.RemoteAddr, pkt.Packet[:pkt.Size])
		return nil
	}
}

// serverUdpHandler  core self handler
func (r *RegServer) serverUdpHandler() node.HandlerFunc {
	return func(ctx context.Context, frame *node.Frame) error {
		logger.Infof("server got packet size:%d, data type: [%v]", frame.Size, util.GetFrameTypeName(util.MsgTypePacket))
		switch frame.FrameType {
		case util.MsgTypeRegisterSuper:
			err := r.register(frame)
			h, err := packet.NewHeader(util.MsgTypeRegisterAck, frame.NetworkId)
			if err != nil {
				logger.Errorf("build resp failed. err: %v", err)
			}
			f, _ := packet.Encode(h)
			frame.Packet = f
			break
		case util.MsgTypePacket:
			p, err := r.cache.GetPeer(frame.UidString(), frame.DstIP.String())
			if err != nil || p == nil {
				return fmt.Errorf("peer %v is not found", frame.DstIP.String())
			}
			logger.Debugf("write packet to peer %v: ", p)
			frame.RemoteAddr = p.GetEndpoint().DstIP()
			frame.Peer = p //change peer to dst peer
			r.PutPktToOutbound(frame)
		case util.MsgTypeQueryPeer:
			peers := r.cache.ListPeers(frame.UidString())
			peerAck := peer.NewPeerPacket()

			for ip, p := range peers {
				info := peer.PeerInfo{
					IP:         net.ParseIP(ip),
					RemoteAddr: *p.GetEndpoint().DstIP(),
				}
				peerAck.Peers = append(peerAck.Peers, info)
			}
			buff, _ := peer.Encode(peerAck)

			newFrame := node.NewFrame()
			copy(newFrame.Packet, buff)
			newFrame.UserId = frame.UserId
			newFrame.RemoteAddr = frame.RemoteAddr
			newFrame.FrameType = util.MsgTypeQueryPeer
			newFrame.Size = len(buff)
			r.PutPktToOutbound(newFrame)
		case util.HandShakeMsgType:
			if _, err := node.CachePeerToLocal(r.key.privateKey, frame, r.cache, nil); err != nil {
				return err
			}
			//build handshake resp
			hpkt := handshake.NewPacket(util.HandShakeMsgTypeAck, frame.UidString())
			hpkt.Header.SrcIP = frame.DstIP
			hpkt.Header.DstIP = frame.SrcIP
			hpkt.PubKey = r.key.privateKey.NewPubicKey()
			buff, err := handshake.Encode(hpkt)
			if err != nil {
				return err
			}

			newFrame := node.NewFrame()
			newFrame.Size = len(buff)
			newFrame.RemoteAddr = frame.RemoteAddr
			copy(newFrame.Packet[:newFrame.Size], buff)
			r.PutPktToOutbound(newFrame)
		}

		return nil
	}
}

func (r *RegServer) register(frame *node.Frame) (err error) {
	p := new(node.Peer)
	ep := nets.NewEndpoint(frame.RemoteAddr.String())
	//ep.SetSrcIP(frame.SrcIP)
	p.SetEndpoint(ep)
	err = r.cache.SetPeer(frame.UidString(), frame.SrcIP.String(), p)
	return
}
