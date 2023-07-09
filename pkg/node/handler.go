package node

import (
	"context"
	. "github.com/topcloudz/fvpn/pkg/handler"
	"github.com/topcloudz/fvpn/pkg/nets"
	"github.com/topcloudz/fvpn/pkg/packet"
	"github.com/topcloudz/fvpn/pkg/packet/handshake"
	"github.com/topcloudz/fvpn/pkg/packet/peer"
	"github.com/topcloudz/fvpn/pkg/packet/register/ack"
	"github.com/topcloudz/fvpn/pkg/security"
	"github.com/topcloudz/fvpn/pkg/util"
)

func (n *Node) tunInHandler() HandlerFunc {
	return func(ctx context.Context, frame *packet.Frame) error {
		//defer frame.Unlock()
		n.PutPktToOutbound(frame)
		return nil
	}
}

// Handle union udp handler
func (n *Node) udpInHandler() HandlerFunc {
	return func(ctx context.Context, frame *packet.Frame) error {
		//dest := ctx.Value("destAddr").(string)
		buff := frame.Packet[:]
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
		case util.MsgTypeQueryPeer:
			logger.Debugf("got list packets response")
			n.handleQueryPeers(frame)
		case util.MsgTypePacket:
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
	hpkt, err := handshake.Decode(frame.Buff)
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

func (n *Node) handleQueryPeers(frame *packet.Frame) {
	peers, _ := peer.Decode(frame.Packet[:])
	logger.Debugf("go peers from remote: %v", peers.Peers)
	for _, info := range peers.Peers {
		ip := info.IP
		addr := info.RemoteAddr

		p := n.NewPeer(info.PubKey)
		p.SetEndpoint(nets.NewEndpoint(addr.String()))
		p.cipher = security.NewCipher(n.privateKey, p.PubKey)
		err := n.cache.SetPeer(frame.UidString(), ip.String(), p)
		if err != nil {
			return
		}
		p.start()
	}
}
