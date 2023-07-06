package node

import (
	"context"
	. "github.com/topcloudz/fvpn/pkg/handler"
	"github.com/topcloudz/fvpn/pkg/nets"
	"github.com/topcloudz/fvpn/pkg/packet"
	"github.com/topcloudz/fvpn/pkg/packet/handshake"
	peerack "github.com/topcloudz/fvpn/pkg/packet/peer/ack"
	"github.com/topcloudz/fvpn/pkg/packet/register/ack"
	"github.com/topcloudz/fvpn/pkg/security"
	"github.com/topcloudz/fvpn/pkg/util"
)

func (n *Node) tunInHandler() HandlerFunc {
	return func(ctx context.Context, frame *packet.Frame) error {

		h, _ := packet.NewHeader(util.MsgTypePacket, "")
		frame.UserId = h.UserId
		headerBuff, err := packet.Encode(h)
		if err != nil {
			return err
		}

		idx := 0
		idx = packet.EncodeBytes(frame.Packet, headerBuff, idx)
		idx = packet.EncodeBytes(frame.Packet, frame.Buff[:frame.Size], idx)

		frame.Size = idx
		frame.FrameType = util.MsgTypePacket
		n.PutPktToOutbound(frame)

		return nil

	}
}

// Handle union udp handler
func (n *Node) udpInHandler() HandlerFunc {
	return func(ctx context.Context, frame *packet.Frame) error {
		//dest := ctx.Value("destAddr").(string)
		buff := frame.Buff[:]
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
			err = CachePeerToLocal(n.privateKey, frame, n.cache)
			if err != nil {
				return err
			}
		}

		return nil
	}
}

func CachePeerToLocal(privateKey security.NoisePrivateKey, frame *packet.Frame, cache CacheFunc) error {
	hpkt, err := handshake.Decode(frame.Packet)
	if err != nil {
		logger.Errorf("invalid handshake packet: %v", err)
		return err
	}

	peer := new(Peer)
	peer.PubKey = hpkt.PubKey
	ep := nets.NewEndpoint(frame.SrcIP.String())
	peer.SetEndpoint(ep)
	err = cache.SetPeer(frame.UidString(), frame.SrcIP.String(), peer)
	peer.cipher = security.NewCipher(privateKey, peer.PubKey)
	peer.start()
	if err != nil {
		return err
	}
	return nil
}
