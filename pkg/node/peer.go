package node

import (
	"github.com/topcloudz/fvpn/pkg/nets"
	"github.com/topcloudz/fvpn/pkg/packet"
	"github.com/topcloudz/fvpn/pkg/packet/handshake"
	"github.com/topcloudz/fvpn/pkg/security"
)

// Peer a destination will have a peer in fvpn, can connect to each other.
// a RegServer also is a peer
type Peer struct {
	PubKey   security.NoisePublicKey
	node     *Node
	endpoint nets.Endpoint //

	queue struct {
		outBound *OutBoundQueue // data to write to dst peer
		inBound  *InBoundQueue  // data write to tun
	}
}

func (p *Peer) start() {
	go p.SendPackets()
	go p.WriteToDevice()
	p.handshake()
}

func (p *Peer) SetEndpoint(ep nets.Endpoint) {
	p.endpoint = ep
}

func (p *Peer) handshake() {
	hpkt := handshake.NewPacket("")
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
	pkt.Lock()
	defer pkt.Unlock()
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
			write, err := p.node.device.Write(pkt.Packet[:pkt.Size])
			if err != nil {
				return
			}

			logger.Debugf("peer %v has write %d packets to device", p, write)
		default:

		}
	}
}
