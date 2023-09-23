package device

import (
	"github.com/tiptopsoft/fvpn/pkg/packet"
	"sync"
	"time"
)

type MemoryPool struct {
	lock sync.Mutex
	pool sync.Pool
	//cond sync.Cond
}

func InitPools() (buffPool *MemoryPool, framePool *MemoryPool) {
	framePool = NewPool(func() any {
		frame := &Frame{
			Ctx:        nil,
			Mutex:      sync.Mutex{},
			Packet:     make([]byte, packet.FvpnPktBuffSize),
			Size:       0,
			NetworkId:  "",
			UserId:     [8]byte{},
			RemoteAddr: nil,
			SrcIP:      nil,
			DstIP:      nil,
			FrameType:  0,
			Peer:       nil,
			Encrypt:    false,
		}
		return frame
	})

	return
}

func NewPool(new func() any) *MemoryPool {
	return &MemoryPool{
		pool: sync.Pool{New: new},
	}
}

func (p *MemoryPool) Get() any {
	//p.lock.Lock()
	return p.pool.Get()
}

func (p *MemoryPool) Put(x any) {
	//defer p.lock.Unlock()
	p.pool.Put(x)
}

func (n *Node) GetBuffer() *[packet.FvpnPktBuffSize]byte {
	return n.pools.buffPool.Get().(*[packet.FvpnPktBuffSize]byte)
}

func (n *Node) PutBuffer(buffPtr *[packet.FvpnPktBuffSize]byte) {
	n.pools.buffPool.Put(buffPtr)
}

func (n *Node) GetFrame() *Frame {
	frame := n.pools.framePool.Get().(*Frame)
	frame.ST = time.Now()
	return frame
}

func (n *Node) PutFrame(framePtr *Frame) {
	framePtr.Size = 0
	framePtr.NetworkId = ""
	framePtr.SrcIP = nil
	framePtr.DstIP = nil
	n.pools.framePool.Put(framePtr)
}
