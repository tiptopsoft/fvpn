package node

import (
	"sync"
)

const (
	QueueOutboundSize       = 1024
	QueueInboundSize        = 1024
	QueueHandshakeBoundSize = 1024
)

type OutBoundQueue struct {
	c  chan *Frame
	wg sync.WaitGroup
}

type InBoundQueue struct {
	c  chan *Frame
	wg sync.WaitGroup
}

type handshakeBound struct {
	c  chan *Frame
	wg sync.WaitGroup
}

func NewOutBoundQueue() *OutBoundQueue {
	q := &OutBoundQueue{
		c: make(chan *Frame, QueueInboundSize),
	}
	q.wg.Add(1)
	go func() {
		q.wg.Wait()
		close(q.c)
	}()

	return q
}

func (o *OutBoundQueue) PutPktToOutbound(pkt *Frame) {
	o.c <- pkt
}

func (o *OutBoundQueue) GetPktFromOutbound() chan *Frame {
	return o.c
}

func (o *InBoundQueue) PutPktToInbound(pkt *Frame) {
	o.c <- pkt
}

func (o *InBoundQueue) GetPktFromInbound() chan *Frame {
	return o.c
}

func NewInBoundQueue() *InBoundQueue {
	q := &InBoundQueue{
		c: make(chan *Frame, QueueOutboundSize),
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
		c: make(chan *Frame, QueueHandshakeBoundSize),
	}
	q.wg.Add(1)
	go func() {
		q.wg.Wait()
		close(q.c)
	}()

	return q
}
