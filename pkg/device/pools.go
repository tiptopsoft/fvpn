package device

import (
	"github.com/tiptopsoft/fvpn/pkg/packet"
	"sync"
)

type PacketPool struct {
	pool sync.Pool
}

func NewPool() *PacketPool {
	return &PacketPool{
		pool: sync.Pool{New: func() interface{} {
			logger.Infof(">>>>>>>>>>>new pool buffer")
			return new([packet.FvpnPktBuffSize]byte)
		}},
	}
}

func (p *PacketPool) Get() *[packet.FvpnPktBuffSize]byte {
	return p.pool.Get().(*[packet.FvpnPktBuffSize]byte)
}

func (p *PacketPool) Put(buffPtr *[packet.FvpnPktBuffSize]byte) {
	p.pool.Put(buffPtr)
}
