package packet

import (
	"github.com/topcloudz/fvpn/pkg/option"
	"sync"
)

type Frame struct {
	sync.Mutex
	Buff      []byte //max length 2000
	Packet    []byte
	NetworkId string
}

func NewFrame() *Frame {
	return &Frame{
		Buff:   make([]byte, option.FVPN_PKT_BUFF_SIZE),
		Packet: make([]byte, option.FVPN_PKT_BUFF_SIZE),
	}
}
