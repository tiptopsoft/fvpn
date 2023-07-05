package packet

import (
	"context"
	"encoding/hex"
	"github.com/topcloudz/fvpn/pkg/security"
	"net"
	"sync"
)

const (
	FvpnPktBuffSize = 2048
)

type Frame struct {
	Ctx context.Context
	sync.Mutex
	Buff       []byte //max length 2000
	Packet     []byte
	Size       int
	NetworkId  string
	UserId     [8]byte
	SrcAddr    *net.UDPAddr
	RemoteAddr string //inner ip
	TargetAddr *net.UDPAddr
	FrameType  uint16
	ciper      security.CipherFunc
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

func (f *Frame) SrcIP() string {
	return f.SrcAddr.IP.To4().String()
}

func (f *Frame) Context() context.Context {
	return f.Ctx
}
