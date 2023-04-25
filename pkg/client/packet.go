package client

import (
	"github.com/topcloudz/fvpn/pkg/option"
	"sync"
)

type Frame struct {
	sync.Mutex
	buff      []byte //max length 2000
	packet    []byte
	networkId string
}

func NewFrame() *Frame {
	return &Frame{
		buff:   make([]byte, option.FVPN_PKT_BUFF_SIZE),
		packet: make([]byte, option.FVPN_PKT_BUFF_SIZE),
	}
}
