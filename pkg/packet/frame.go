package packet

import (
	"errors"
	"github.com/topcloudz/fvpn/pkg/cache"
	"github.com/topcloudz/fvpn/pkg/option"
	"github.com/topcloudz/fvpn/pkg/util"
	"sync"
)

type Frame struct {
	sync.Mutex
	Buff      []byte //max length 2000
	Packet    []byte
	Size      int
	NetworkId string
}

func NewFrame() *Frame {
	return &Frame{
		Buff:   make([]byte, option.FVPN_PKT_BUFF_SIZE),
		Packet: make([]byte, option.FVPN_PKT_BUFF_SIZE),
	}
}

func (f *Frame) GetNodeInfo(cache cache.Cache) (*cache.NodeInfo, error) {
	destMac, err := util.GetMacAddr(f.Packet)
	if err != nil {
		return nil, errors.New("no data exists")
	}
	return cache.GetNodeInfo(destMac)
}
