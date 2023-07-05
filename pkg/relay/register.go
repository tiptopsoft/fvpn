package relay

import (
	"github.com/topcloudz/fvpn/pkg/cache"
	"net"
)

func (r *RegServer) registerAck(peerAddr *net.UDPAddr, srcMac net.HardwareAddr, srcIP net.IP, networkId string) error {

	ackNode := &cache.Endpoint{
		Addr:      peerAddr,
		NetworkId: "",
		MacAddr:   srcMac,
		IP:        srcIP,
		Port:      0,
	}

	ackNode.NatIP = peerAddr.IP
	ackNode.NatPort = uint16(peerAddr.Port)

	r.cache.SetCache(networkId, srcIP.String(), ackNode)
	logger.Debugf("node register success, networkId: %s, ip: %v, natIP: %v, natPort: %d", networkId, srcIP.String(), peerAddr.IP, peerAddr.Port)
	return nil
}
