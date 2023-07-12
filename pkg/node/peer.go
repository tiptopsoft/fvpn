package node

import (
	"github.com/topcloudz/fvpn/pkg/nets"
	"github.com/topcloudz/fvpn/pkg/packet"
	"github.com/topcloudz/fvpn/pkg/packet/handshake"
	"github.com/topcloudz/fvpn/pkg/security"
	"github.com/topcloudz/fvpn/pkg/util"
	"net"
	"sync"
	"time"
)

// Peer a destination will have a peer in fvpn, can connect to each other.
// a RegServer also is a peer
type Peer struct {
	p2p      bool
	lock     sync.Mutex
	status   bool
	PubKey   security.NoisePublicKey
	node     *Node
	endpoint nets.Endpoint //

	queue struct {
		outBound *OutBoundQueue // data to write to dst peer
		inBound  *InBoundQueue  // data write to tun
	}
	cipher security.CipherFunc
}

func (p *Peer) GetCodec() security.CipherFunc {
	return p.cipher
}

func (p *Peer) start() {
	p.lock.Lock()
	defer p.lock.Unlock()
	if p.status == true {
		logger.Debugf("peer has started: %v", p.endpoint.DstToString())
		return
	}

	p.status = true
	if p.node != nil && p.node.mode == 1 {
		p.PubKey = p.node.privateKey.NewPubicKey() //use to send to remote for exchange pubKey
		go p.SendPackets()
		//go p.WriteToDevice()

		go func() {
			timer := time.NewTimer(time.Second * 10)
			defer timer.Stop()
			for {
				select {
				case <-timer.C:
					logger.Debugf("sending keepalive....")
					p.keepalive()
					timer.Reset(time.Second * 10)
				}
			}
		}()
	}
}

func (p *Peer) SetEndpoint(ep nets.Endpoint) {
	p.endpoint = ep
}

func (p *Peer) GetEndpoint() nets.Endpoint {
	return p.endpoint
}

func (p *Peer) handshake(dstIP net.IP) {
	hpkt := handshake.NewPacket(util.HandShakeMsgType, util.UCTL.UserId)
	hpkt.Header.SrcIP = p.node.device.Addr()
	hpkt.Header.DstIP = dstIP
	hpkt.PubKey = p.PubKey
	buff, err := handshake.Encode(hpkt)
	if err != nil {
		return
	}

	size := len(buff)
	f := NewFrame()
	copy(f.Packet[:size], buff)
	f.Size = size
	//cache a peer
	//ep := nets.NewEndpoint(p.endpoint.DstToString())
	//p.SetEndpoint(ep)
	//err = p.node.cache.SetPeer(handler.UCTL.UserId, p.endpoint.DstToString(), p)
	//if err != nil {
	//	logger.Error("init cache peer failed.")
	//	return
	//}
	logger.Debugf("sending pubkey to: %v, pubKey: %v", dstIP.String(), p.PubKey)
	p.PutPktToOutbound(f)
}

func (p *Peer) PutPktToOutbound(pkt *Frame) {
	//pkt.Lock()
	//defer pkt.Unlock()
	p.queue.outBound.c <- pkt
}

func (p *Peer) SendPackets() {
	for {
		select {
		case pkt := <-p.queue.outBound.c:
			send, err := p.node.net.bind.Send(pkt.Packet[:pkt.Size], p.endpoint)
			if err != nil {
				logger.Error(err)
				continue
			}
			logger.Debugf("peer %v has send %d packets to %s, buff: %v", p, send, p.endpoint.DstToString(), pkt.Packet[:pkt.Size])
		default:

		}
	}
}

func (p *Peer) keepalive() {
	pkt, err := packet.NewHeader(util.KeepaliveMsgType, "")
	if err != nil {
		return
	}
	buff, err := packet.Encode(pkt)
	if err != nil {
		return
	}
	size := len(buff)
	f := NewFrame()
	copy(f.Packet[:size], buff)
	f.Size = size

	p.PutPktToOutbound(f)
}
