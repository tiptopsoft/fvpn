package server

import (
	"github.com/topcloudz/fvpn/pkg/cache"
	"golang.org/x/sys/unix"
	"net"
)

func (r *RegServer) registerAck(peerAddr *unix.SockaddrInet4, srcMac net.HardwareAddr, srcIP net.IP, networkId string) error {

	ackNode := &cache.NodeInfo{
		Socket:    r.socket,
		Addr:      peerAddr,
		NetworkId: "",
		MacAddr:   srcMac,
		IP:        srcIP,
		Port:      0,
	}

	r.cache.SetCache(networkId, srcIP.String(), ackNode)
	logger.Debugf("node register success, networkId: %s, ip: %v, natIP: %v, natPort: %d", networkId, srcIP.String(), peerAddr.Addr, peerAddr.Port)
	return nil
}
