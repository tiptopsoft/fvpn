package device

import (
	"context"
	"github.com/topcloudz/fvpn/pkg/util"
	"net"

	"github.com/topcloudz/fvpn/pkg/addr"
	"github.com/topcloudz/fvpn/pkg/cache"
	"github.com/topcloudz/fvpn/pkg/handler"
	"github.com/topcloudz/fvpn/pkg/log"
	"github.com/topcloudz/fvpn/pkg/packet"
	"github.com/topcloudz/fvpn/pkg/packet/forward"
	"github.com/topcloudz/fvpn/pkg/socket"
	"github.com/topcloudz/fvpn/pkg/tuntap"
)

var (
	logger = log.Log()
)

type DeviceHandler struct {
	net    socket.Interface
	cache  cache.PeersCache
	device *tuntap.Tuntap
}

func New(device *tuntap.Tuntap, netSocket socket.Interface, cache cache.PeersCache) handler.Handler {
	return DeviceHandler{
		net:    netSocket,
		device: device,
		cache:  cache,
	}
}

func (dh DeviceHandler) Handle(ctx context.Context, buff []byte) error {
	device := dh.device
	var err error
	destMac := util.GetMacAddr(buff)

	broad := addr.IsBroadCast(destMac)
	//broad frame, go through supernode
	fp := forward.NewPacket()
	fp.SrcMac, err = addr.GetMacAddrByDev(device.Name)
	if err != nil {
		logger.Errorf("get src mac failed, err: %v", err)
	}
	fp.DstMac, err = net.ParseMAC(destMac)
	if err != nil {
		logger.Errorf("get src mac failed, err: %v", err)
	}

	bs, err := fp.Encode()
	if err != nil {
		logger.Errorf("encode forward failed, err: %v", err)
	}

	idx := 0
	newPacket := make([]byte, 2048)
	idx = packet.EncodeBytes(newPacket, bs, idx)
	idx = packet.EncodeBytes(newPacket, buff[:], idx)
	if broad {
		dh.write2Net(newPacket[:idx])
	} else {
		// go p2p
		logger.Infof("find peer in client, destMac: %v", destMac)
		p := cache.FindPeer(dh.cache, destMac)
		if p == nil {
			dh.write2Net(newPacket[:idx])
			logger.Warnf("peer not found, go through super cache")
		} else {
			dh.write2Net(newPacket[:idx])
		}
	}
	return nil
}

// use host socket write to destination, superNode or use p2p
func (dh DeviceHandler) write2Net(b []byte) {
	logger.Debugf("tap write to net packet: (%v)", b)
	if _, err := dh.net.Write(b); err != nil {
		logger.Errorf("tap write to net failed. (%v)", err)
	}
}
