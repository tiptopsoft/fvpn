package node

import (
	"github.com/topcloudz/fvpn/pkg/packet"
	"sync"
)

const (
	QueueOutboundSize       = 1024
	QueueInboundSize        = 1024
	QueueHandshakeBoundSize = 1024
)

type OutBoundQueue struct {
	c  chan *packet.Frame
	wg sync.WaitGroup
}

type InBoundQueue struct {
	c  chan *packet.Frame
	wg sync.WaitGroup
}

type handshakeBound struct {
	c  chan *packet.Frame
	wg sync.WaitGroup
}

func NewOutBoundQueue() *OutBoundQueue {
	q := &OutBoundQueue{
		c: make(chan *packet.Frame, QueueInboundSize),
	}
	q.wg.Add(1)
	go func() {
		q.wg.Wait()
		close(q.c)
	}()

	return q
}

func (o *OutBoundQueue) PutPktToOutbound(pkt *packet.Frame) {
	pkt.Lock()
	defer pkt.Unlock()
	o.c <- pkt
}

func (o *OutBoundQueue) GetPktFromOutbound() chan *packet.Frame {
	return o.c
}

func (o *InBoundQueue) PutPktToInbound(pkt *packet.Frame) {
	pkt.Lock()
	defer pkt.Unlock()
	o.c <- pkt
}

func (o *InBoundQueue) GetPktFromInbound() chan *packet.Frame {
	return o.c
}

func NewInBoundQueue() *InBoundQueue {
	q := &InBoundQueue{
		c: make(chan *packet.Frame, QueueOutboundSize),
	}
	q.wg.Add(1)
	go func() {
		q.wg.Wait()
		close(q.c)
	}()

	return q
}

func newHandshakeQueue() *handshakeBound {
	q := &handshakeBound{
		c: make(chan *packet.Frame, QueueHandshakeBoundSize),
	}
	q.wg.Add(1)
	go func() {
		q.wg.Wait()
		close(q.c)
	}()

	return q
}
