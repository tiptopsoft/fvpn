package relay

import (
	"github.com/tiptopsoft/fvpn/pkg/device"
	"github.com/tiptopsoft/fvpn/pkg/packet"
)

func (n *RegServer) GetBuffer() *[packet.FvpnPktBuffSize]byte {
	return n.pools.buffPool.Get().(*[packet.FvpnPktBuffSize]byte)
}

func (n *RegServer) PutBuffer(buffPtr *[packet.FvpnPktBuffSize]byte) {
	n.pools.buffPool.Put(buffPtr)
}

func (n *RegServer) GetFrame() *device.Frame {
	return n.pools.framePool.Get().(*device.Frame)
}

func (n *RegServer) PutFrame(framePtr *device.Frame) {
	n.pools.framePool.Put(framePtr)
}
