package packet

import (
	"net"
	"sync"
)

const (
	FvpnPktBuffSize = 2048
)

type Frame struct {
	sync.Mutex
	Buff      []byte //max length 2000
	Packet    []byte
	Size      int
	NetworkId string
	UserId    []byte
	//PubKey      string
	SrcAddr    *net.UDPAddr
	RemoteAddr string //inner ip
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
