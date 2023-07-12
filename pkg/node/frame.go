package node

import (
	"context"
	"encoding/hex"
	"github.com/topcloudz/fvpn/pkg/packet"
	"net"
	"sync"
)

type Frame struct {
	Ctx context.Context
	sync.Mutex
	Buff       []byte
	Packet     []byte
	Size       int
	NetworkId  string
	UserId     [8]byte
	RemoteAddr *net.UDPAddr //remote socket endpoint
	SrcIP      net.IP
	DstIP      net.IP
	FrameType  uint16
	Peer       *Peer
}

func (f *Frame) GetPeer() *Peer {
	return f.Peer
}

func NewFrame() *Frame {
	return &Frame{
		Ctx:    context.Background(),
		Buff:   make([]byte, packet.FvpnPktBuffSize),
		Packet: make([]byte, packet.FvpnPktBuffSize),
	}
}

func (f *Frame) Clear() {
	buf := make([]byte, packet.FvpnPktBuffSize)
	copy(f.Packet, buf)
}

func (f *Frame) UidString() string {
	return hex.EncodeToString(f.UserId[:])
}

func (f *Frame) Context() context.Context {
	return f.Ctx
}
