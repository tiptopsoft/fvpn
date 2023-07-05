package relay

import (
	"context"
	"encoding/hex"
	"github.com/topcloudz/fvpn/pkg/handler"
	"github.com/topcloudz/fvpn/pkg/nets"
	"github.com/topcloudz/fvpn/pkg/node"
	"github.com/topcloudz/fvpn/pkg/packet"
	"github.com/topcloudz/fvpn/pkg/packet/handshake"
	handack "github.com/topcloudz/fvpn/pkg/packet/handshake/ack"
	"github.com/topcloudz/fvpn/pkg/packet/header"
	"github.com/topcloudz/fvpn/pkg/util"
)

func (r *RegServer) ReadFromUdp() {
	logger.Infof("start a udp loop")
	for {
		ctx := context.Background()
		frame := packet.NewFrame()
		frame.Ctx = ctx
		n, addr, err := r.conn.ReadFromUDP(frame.Buff[:])
		frame.Packet = frame.Buff[:n]
		logger.Debugf("Read from udp %d byte, data: %v", n, frame.Packet)

		packetHeader, err := util.GetPacketHeader(frame.Packet[:])
		if err != nil {
			logger.Errorf("get header falied. %v", err)
			continue
		}
		networkId := hex.EncodeToString(packetHeader.NetworkId[:])
		frame.Size = n
		frame.FrameType = packetHeader.Flags
		frame.SrcAddr = addr
		//frame.PubKey = hex.EncodeToString(packetHeader.PubKey[:])
		frame.NetworkId = networkId
		frame.UserId = packetHeader.UserId
		if err != nil || n < 0 {
			logger.Warnf("no data exists")
			continue
		}

		r.PutPktToInbound(frame)
	}
}

func (r *RegServer) writeUdpHandler() handler.HandlerFunc {
	return func(ctx context.Context, pkt *packet.Frame) error {
		n, err := r.conn.WriteToUDP(pkt.Packet[:pkt.Size], pkt.TargetAddr)
		if err != nil {
			return err
		}
		logger.Debugf("registry write %d size to %v", n, pkt.TargetAddr)
		return nil
	}
}

// serverUdpHandler  core self handler
func (r *RegServer) serverUdpHandler() handler.HandlerFunc {
	return func(ctx context.Context, frame *packet.Frame) error {
		defer frame.Unlock()
		data := frame.Packet[:frame.Size]
		switch frame.FrameType {
		case util.MsgTypeRegisterSuper:
			err := r.register(frame)
			h, err := header.NewHeader(util.MsgTypeRegisterAck, frame.NetworkId)
			if err != nil {
				logger.Errorf("build resp failed. err: %v", err)
			}
			f, _ := header.Encode(h)
			frame.Packet = f
			frame.TargetAddr = frame.SrcAddr
			break
		case util.MsgTypePacket:
			logger.Infof("server got forward packet size:%d, data: %v", frame.Size, data)
		case util.HandShakeMsgType:
			handPkt, err := handshake.Decode(frame.Packet)
			if err != nil {
				logger.Errorf("invalid handshake packet: %v", err)
				return err
			}

			peer := new(node.Peer)
			peer.PubKey = handPkt.PubKey
			ep := nets.NewEndpoint(frame.SrcIP())
			peer.SetEndpoint(ep)
			err = r.cache.SetPeer(frame.UidString(), frame.SrcIP(), peer)
			if err != nil {
				return err
			}

			//build handshake ack

			hPktack := handack.NewPacket()
			hPktack.PubKey = r.key.privateKey.NewPubicKey()
			buff, err := handack.Encode(hPktack)
			if err != nil {
				return err
			}

			f := packet.NewFrame()
			frame.Size = len(buff)
			frame.TargetAddr = frame.SrcAddr
			copy(f.Packet[:frame.Size], buff)
			r.PutPktToOutbound(f)
		}

		return nil
	}
}

func (r *RegServer) register(frame *packet.Frame) (err error) {
	p := new(node.Peer)
	err = r.cache.SetPeer(frame.UidString(), frame.SrcIP(), p)
	return
}
