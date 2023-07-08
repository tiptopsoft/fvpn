package node

import (
	"github.com/topcloudz/fvpn/pkg/handler"
	"github.com/topcloudz/fvpn/pkg/nets"
	"github.com/topcloudz/fvpn/pkg/packet"
	"github.com/topcloudz/fvpn/pkg/packet/handshake"
	"github.com/topcloudz/fvpn/pkg/security"
	"github.com/topcloudz/fvpn/pkg/util"
	"sync"
	"time"
)

// Peer a destination will have a peer in fvpn, can connect to each other.
// a RegServer also is a peer
type Peer struct {
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
		return
	}
	p.status = true
	if p.node != nil && p.node.mode == 1 {
		go p.SendPackets()
		go p.WriteToDevice()
		//cache peer
		p.handshake()

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

func (p *Peer) handshake() {
	hpkt := handshake.NewPacket(util.HandShakeMsgType, handler.UCTL.UserId)
	hpkt.Header.SrcIP = p.node.device.Addr()
	hpkt.Header.DstIP = p.endpoint.DstIP().IP
	hpkt.PubKey = p.node.pubKey
	buff, err := handshake.Encode(hpkt)
	if err != nil {
		return
	}

	size := len(buff)
	f := packet.NewFrame()
	copy(f.Packet[:size], buff)
	f.Size = size

	p.PutPktToOutbound(f)

	//cache

}

func (p *Peer) PutPktToOutbound(pkt *packet.Frame) {
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

func (p *Peer) WriteToDevice() {
	for {
		select {
		case pkt := <-p.queue.inBound.c:
			write, err := p.node.device.Write(pkt.Packet[packet.HeaderBuffSize:pkt.Size])
			if err != nil {
				return
			}

			logger.Debugf("peer %v has write %d packets to device", p, write)
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
	f := packet.NewFrame()
	copy(f.Packet[:size], buff)
	f.Size = size

	p.PutPktToOutbound(f)
}
