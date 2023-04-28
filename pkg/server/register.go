package server

import (
	"github.com/topcloudz/fvpn/pkg/cache"
	"golang.org/x/sys/unix"
	"net"
)

func (r *RegServer) registerAck(peerAddr unix.Sockaddr, srcMac net.HardwareAddr, srcIP net.IP, networkId string) error {

	ackNode := &cache.NodeInfo{
		Socket:    r.socket,
		Addr:      peerAddr,
		NetworkId: "",
		MacAddr:   srcMac,
		IP:        srcIP,
		Port:      0,
	}

	r.cache.SetCache(networkId, srcIP.String(), ackNode)
	return nil
}
