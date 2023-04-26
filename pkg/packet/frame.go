package packet

import (
	"github.com/topcloudz/fvpn/pkg/cache"
	"github.com/topcloudz/fvpn/pkg/option"
	"github.com/topcloudz/fvpn/pkg/util"
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

func (f *Frame) GetNodeInfo(cache cache.Cache) (*cache.NodeInfo, error) {
	destMac := util.GetMacAddr(f.Packet)
	return cache.GetNodeInfo(destMac)
}
