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

type outBoundQueue struct {
	c  chan *packet.Frame
	wg sync.WaitGroup
}

type inBoundQueue struct {
	c  chan *packet.Frame
	wg sync.WaitGroup
}

type handshakeBound struct {
	c  chan *packet.Frame
	wg sync.WaitGroup
}

func newOutBoundQueue() *outBoundQueue {
	q := &outBoundQueue{
		c: make(chan *packet.Frame, QueueInboundSize),
	}
	q.wg.Add(1)
	go func() {
		q.wg.Wait()
		close(q.c)
	}()

	return q
}

func newInBoundQueue() *inBoundQueue {
	q := &inBoundQueue{
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
