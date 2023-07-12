package relay

import (
	"context"
	"fmt"
	"github.com/topcloudz/fvpn/pkg/handler"
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
		frame := packet.NewFrame()
		frame.Ctx = ctx
		n, addr, err := r.conn.ReadFromUDP(frame.Buff[:])
		if err != nil || n < 0 {
			logger.Error("no data exists")
			continue
		}
		logger.Debugf("Read from udp %d byte, data: %v", n, frame.Buff[:n])
		copy(frame.Packet, frame.Buff)
		packetHeader, err := util.GetPacketHeader(frame.Buff[:])
		if err != nil {
			logger.Errorf("get header falied. %v", err)
			continue
		}
		frame.Size = n
		frame.FrameType = packetHeader.Flags
		frame.RemoteAddr = addr
		frame.SrcIP = packetHeader.SrcIP
		frame.DstIP = packetHeader.DstIP
		frame.UserId = packetHeader.UserId

		r.PutPktToInbound(frame)
	}
}

func (r *RegServer) writeUdpHandler() handler.HandlerFunc {
	return func(ctx context.Context, pkt *packet.Frame) error {
		n, err := r.conn.WriteToUDP(pkt.Packet[:pkt.Size], pkt.RemoteAddr)
		if err != nil {
			return err
		}
		logger.Debugf("registry write %d size to %v, data: %v", n, pkt.RemoteAddr, pkt.Packet[:pkt.Size])
		return nil
	}
}

// serverUdpHandler  core self handler
func (r *RegServer) serverUdpHandler() handler.HandlerFunc {
	return func(ctx context.Context, frame *packet.Frame) error {
		data := frame.Packet[:frame.Size]
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
			logger.Infof("server got forward packet size:%d, data: %v", frame.Size, data)
			peer, err := r.cache.GetPeer(frame.UidString(), frame.DstIP.String())
			if err != nil || peer == nil {
				return fmt.Errorf("peer %v is not found", frame.DstIP.String())
			}

			logger.Debugf("write packet to peer %v: ", peer)

			frame.RemoteAddr = peer.GetEndpoint().DstIP()
			r.PutPktToOutbound(frame)
		case util.MsgTypeQueryPeer:
			logger.Debug("server got list peers packet")
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

			newFrame := packet.NewFrame()
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

			newFrame := packet.NewFrame()
			newFrame.Size = len(buff)
			newFrame.RemoteAddr = frame.RemoteAddr
			copy(newFrame.Packet[:newFrame.Size], buff)
			r.PutPktToOutbound(newFrame)
		}

		return nil
	}
}

func (r *RegServer) register(frame *packet.Frame) (err error) {
	p := new(node.Peer)
	ep := nets.NewEndpoint(frame.RemoteAddr.String())
	p.SetEndpoint(ep)
	err = r.cache.SetPeer(frame.UidString(), frame.SrcIP.String(), p)
	return
}
