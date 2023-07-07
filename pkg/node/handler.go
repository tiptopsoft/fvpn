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
		defer frame.Unlock()
		h, _ := packet.NewHeader(util.MsgTypePacket, UCTL.UserId)
		frame.UserId = h.UserId
		h.SrcIP = frame.SrcIP
		h.DstIP = frame.DstIP
		headerBuff, err := packet.Encode(h)
		if err != nil {
			return err
		}

		newFrame := packet.NewFrame()
		size := packet.HeaderBuffSize + frame.Size
		newFrame.Size = size
		newFrame.SrcIP = frame.SrcIP
		newFrame.DstIP = frame.DstIP
		newFrame.UserId = frame.UserId

		idx := 0
		idx = packet.EncodeBytes(newFrame.Packet, headerBuff, idx)
		idx = packet.EncodeBytes(newFrame.Packet, frame.Packet[:frame.Size], idx)

		newFrame.FrameType = util.MsgTypePacket
		n.PutPktToOutbound(newFrame)

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
			n.PutPktToInbound(frame)
		case util.HandShakeMsgType:
			//cache dst peer when receive a handshake
			err = CachePeerToLocal(n.privateKey, frame, n.cache)
			if err != nil {
				return err
			}
			//build handshake resp
			hPktack := handshake.NewPacket(util.HandShakeMsgTypeAck, frame.UidString())
			hPktack.Header.SrcIP = frame.DstIP
			hPktack.Header.DstIP = frame.SrcIP
			hPktack.PubKey = n.privateKey.NewPubicKey()
			buff, err := handshake.Encode(hPktack)
			if err != nil {
				return err
			}

			frame.Packet = buff
			frame.Size = len(buff)
		case util.HandShakeMsgTypeAck:
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

	peer, err := cache.GetPeer(UCTL.UserId, frame.SrcIP.String())
	if err != nil || peer == nil {
		peer = new(Peer)
	}

	peer.PubKey = hpkt.PubKey
	ep := nets.NewEndpoint(frame.RemoteAddr.String())
	peer.SetEndpoint(ep)
	peer.cipher = security.NewCipher(privateKey, peer.PubKey)
	err = cache.SetPeer(frame.UidString(), frame.SrcIP.String(), peer)
	peer.start()

	if err != nil {
		return err
	}
	return nil
}
