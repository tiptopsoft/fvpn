package node

import (
	"context"
	"encoding/hex"
	. "github.com/topcloudz/fvpn/pkg/handler"
	"github.com/topcloudz/fvpn/pkg/nets"
	"github.com/topcloudz/fvpn/pkg/packet"
	"github.com/topcloudz/fvpn/pkg/packet/handshake"
	"github.com/topcloudz/fvpn/pkg/packet/header"
	peerack "github.com/topcloudz/fvpn/pkg/packet/peer/ack"
	"github.com/topcloudz/fvpn/pkg/packet/register/ack"
	"github.com/topcloudz/fvpn/pkg/util"
)

func (n *Node) tunInHandler() HandlerFunc {
	return func(ctx context.Context, frame *packet.Frame) error {
		networkId := ctx.Value("networkId").(string)
		h, _ := header.NewHeader(util.MsgTypePacket, networkId)
		frame.UserId = h.UserId
		headerBuff, err := header.Encode(h)
		if err != nil {
			return err
		}

		idx := 0
		newPacket := make([]byte, 2048)
		idx = packet.EncodeBytes(newPacket, headerBuff, idx)
		idx = packet.EncodeBytes(newPacket, frame.Buff[:frame.Size], idx)

		frame.Packet = newPacket[:idx]
		frame.FrameType = util.MsgTypePacket

		return nil

	}
}

// Handle union udp handler
func (n *Node) udpInHandler() HandlerFunc {
	return func(ctx context.Context, frame *packet.Frame) error {
		//dest := ctx.Value("destAddr").(string)
		buff := frame.Buff[:]
		headerBuff, err := header.Decode(buff)
		if err != nil {
			return err
		}

		frame.NetworkId = hex.EncodeToString(headerBuff.NetworkId[:])
		frame.RemoteAddr = hex.EncodeToString(headerBuff.UserId[:])

		//frame.FrameType = headerBuff.Flags
		switch headerBuff.Flags {
		case util.MsgTypeRegisterAck:
			regAck, err := ack.Decode(buff)
			if err != nil {
				return err
			}
			logger.Debugf("register success, got server server ack: (%v)", regAck)
			break
		case util.MsgTypeQueryPeer:
			logger.Debugf("start get query response")
			peerPacketAck, err := peerack.Decode(buff)
			if err != nil {
				return err
			}
			infos := peerPacketAck.NodeInfos
			logger.Infof("got server peers: (%v)", infos)

			break
		case util.MsgTypePacket:
			frame.Packet = buff[:]
		case util.HandShakeMsgType:
			//cache dst peer when receive a handshake
			err = CachePeerToLocal(frame, n.cache)
			if err != nil {
				return err
			}
		}

		return nil
	}
}

func CachePeerToLocal(frame *packet.Frame, cache CacheFunc) error {
	hpkt, err := handshake.Decode(frame.Packet)
	if err != nil {
		logger.Errorf("invalid handshake packet: %v", err)
		return err
	}

	peer := new(Peer)
	peer.PubKey = hpkt.PubKey
	ep := nets.NewEndpoint(frame.SrcIP())
	peer.SetEndpoint(ep)
	err = cache.SetPeer(frame.UidString(), hpkt.SrcIP.String(), peer)
	if err != nil {
		return err
	}
	return nil
}
