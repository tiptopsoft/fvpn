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
			logger.Debugf("got handshake msg type, data: %v", frame.Packet[:frame.Size])
			_, err := CachePeerToLocal(n.privateKey, frame, n.cache, n)
			if err != nil {
				return err
			}
			//build handshake resp
			//hPktack := handshake.NewPacket(util.HandShakeMsgTypeAck, frame.UidString())
			//logger.Debugf("got packet srcIP: %v, dstIP: %v", frame.SrcIP, frame.DstIP)
			//hPktack.Header.SrcIP = frame.DstIP //dstIP = 2
			//hPktack.Header.DstIP = frame.SrcIP //srcIP = 1
			//hPktack.PubKey = n.privateKey.NewPubicKey()
			//buff, err := handshake.Encode(hPktack)
			//if err != nil {
			//	return err
			//}
			//
			//frame.Packet = buff
			//frame.Size = len(buff)
			//frame.DstIP = frame.SrcIP //dstIP = 1
			//n.PutPktToOutbound(frame)
		case util.HandShakeMsgTypeAck: //use for relay
			//cache dst peer when receive a handshake
			logger.Debugf("got handshake msg type in handshake ack, data: %v", frame.Packet[:frame.Size])
			_, err = CachePeerToLocal(n.privateKey, frame, n.cache, n)
			if err != nil {
				return err
			}
		case util.KeepaliveMsgType:
			logger.Debugf("got keepalived packets from :%v, data: %v", frame.RemoteAddr, frame.Packet[:frame.Size])
		}

		return nil
	}
}

func CachePeerToLocal(privateKey security.NoisePrivateKey, frame *packet.Frame, cache CacheFunc, node *Node) (*Peer, error) {
	hpkt, err := handshake.Decode(frame.Buff)
	if err != nil {
		logger.Errorf("invalid handshake packet: %v", err)
		return nil, err
	}

	logger.Debugf("got remote peer: %v, pubKey: %v", frame.SrcIP.String(), hpkt.PubKey)
	p, err := cache.GetPeer(UCTL.UserId, frame.SrcIP.String())
	if err != nil || p == nil {
		p = node.NewPeer(hpkt.PubKey)
	}
	p.node = node
	ep := nets.NewEndpoint(frame.RemoteAddr.String())
	p.SetEndpoint(ep)
	p.cipher = security.NewCipher(privateKey, hpkt.PubKey)
	p.p2p = true
	err = cache.SetPeer(frame.UidString(), frame.SrcIP.String(), p)
	p.start()

	if err != nil {
		return nil, err
	}

	return p, nil
}

func (n *Node) handleQueryPeers(frame *packet.Frame) {
	peers, _ := peer.Decode(frame.Packet[:])
	logger.Debugf("go peers from remote: %v", peers.Peers)
	for _, info := range peers.Peers {
		ip := info.IP
		if ip.String() == n.device.IPToString() {
			continue
		}
		addr := info.RemoteAddr
		p, err := n.cache.GetPeer(frame.UidString(), ip.String())
		if err != nil || p == nil {
			p = n.NewPeer(security.NoisePublicKey{}) //now has no pubKey
			p.SetEndpoint(nets.NewEndpoint(addr.String()))
			err = n.cache.SetPeer(frame.UidString(), ip.String(), p)
		} else {
			if p.endpoint.DstToString() != addr.String() {
				p.SetEndpoint(nets.NewEndpoint(addr.String()))
			}
		}

		if err != nil {
			return
		}
		//p.start()
		p.start()
		p.handshake(ip)
	}
}
