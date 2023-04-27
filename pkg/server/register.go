package server

import (
	"github.com/topcloudz/fvpn/pkg/addr"
	"github.com/topcloudz/fvpn/pkg/packet/register/ack"
	"golang.org/x/sys/unix"
	"net"
)

func registerAck(peerAddr unix.Sockaddr, srcMac net.HardwareAddr) ([]byte, error) {
	endpoint, err := addr.New(srcMac)
	if err != nil {
		return nil, err
	}
	p := ack.NewPacket()
	p.RegMac = endpoint.Mac
	p.AutoIP = endpoint.IP
	p.Mask = endpoint.Mask

	//ackNode := &cache.NodeInfo{
	//	Socket:    r.socket,
	//	Addr:      peerAddr,
	//	NetworkId: "",
	//	MacAddr:   endpoint.Mac,
	//	IP:        endpoint.IP,
	//	Port:      0,
	//}

	//r.cache.SetCache(endpoint.Mac.String(), ackNode)
	//r.cache.Nodes[endpoint.Mac.String()] = ackNode
	//r.cache.IPNodes[endpoint.IP.String()] = ackNode
	return p.Encode()
}
