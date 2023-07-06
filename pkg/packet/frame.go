package packet

import (
	"context"
	"encoding/hex"
	"net"
	"sync"
)

const (
	FvpnPktBuffSize = 2048
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
}

func NewFrame() *Frame {
	return &Frame{
		Buff:   make([]byte, FvpnPktBuffSize),
		Packet: make([]byte, FvpnPktBuffSize),
	}
}

func (f *Frame) Clear() {
	buf := make([]byte, FvpnPktBuffSize)
	copy(f.Packet, buf)
}

func (f *Frame) UidString() string {
	return hex.EncodeToString(f.UserId[:])
}

func (f *Frame) Context() context.Context {
	return f.Ctx
}
