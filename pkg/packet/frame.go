package packet

import (
	"github.com/topcloudz/fvpn/pkg/cache"
	"github.com/topcloudz/fvpn/pkg/option"
	"golang.org/x/sys/unix"
	"sync"
)

type Frame struct {
	sync.Mutex
	Buff       []byte //max length 2000
	Packet     []byte
	Size       int
	NetworkId  string
	SrcAddr    unix.Sockaddr
	RemoteAddr string //inner ip
	FrameType  uint16
	Type       uint16
	Self       *cache.Endpoint
	Target     *cache.Endpoint
}

func NewFrame() *Frame {
	return &Frame{
		Buff:   make([]byte, option.FVPN_PKT_BUFF_SIZE),
		Packet: make([]byte, option.FVPN_PKT_BUFF_SIZE),
	}
}
