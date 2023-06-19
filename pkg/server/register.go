package server

import (
	"fmt"
	"github.com/topcloudz/fvpn/pkg/cache"
	"golang.org/x/sys/unix"
	"net"
)

func (r *RegServer) registerAck(peerAddr *unix.SockaddrInet4, srcMac net.HardwareAddr, srcIP net.IP, networkId string) error {

	ackNode := &cache.Endpoint{
		Addr:      peerAddr,
		NetworkId: "",
		MacAddr:   srcMac,
		IP:        srcIP,
		Port:      0,
	}

	ackNode.NatIP = net.ParseIP(fmt.Sprintf("%d.%d.%d.%d", peerAddr.Addr[0], peerAddr.Addr[1], peerAddr.Addr[2], peerAddr.Addr[3]))
	ackNode.NatPort = uint16(peerAddr.Port)

	r.cache.SetCache(networkId, srcIP.String(), ackNode)
	logger.Debugf("node register success, networkId: %s, ip: %v, natIP: %v, natPort: %d", networkId, srcIP.String(), peerAddr.Addr, peerAddr.Port)
	return nil
}
