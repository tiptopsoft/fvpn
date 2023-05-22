package server

import (
	"fmt"
	"github.com/topcloudz/fvpn/pkg/cache"
	"github.com/topcloudz/fvpn/pkg/packet/peer/ack"
	"golang.org/x/sys/unix"
	"net"
)

func getPeerInfo(nodes []*cache.NodeInfo) ([]ack.EdgeInfo, uint8, error) {
	var result []ack.EdgeInfo

	for _, peer := range nodes {

		nat := peer.Addr.(*unix.SockaddrInet4)
		addr := nat.Addr
		port := nat.Port
		natIp := net.ParseIP(fmt.Sprintf("%d.%d.%d.%d", addr[0], addr[1], addr[2], addr[3]))
		natPort := uint16(port)
		info := ack.EdgeInfo{
			Mac:     peer.MacAddr,
			IP:      peer.IP,
			Port:    peer.Port,
			NatIp:   natIp,
			NatPort: natPort,
		}
		result = append(result, info)
	}

	return result, uint8(len(result)), nil
}

func peerAckBuild(infos []ack.EdgeInfo, size uint8, networkId string) ([]byte, error) {
	peerPacket := ack.NewPacket(networkId)
	peerPacket.Size = size
	peerPacket.NodeInfos = infos

	return ack.Encode(peerPacket)
}
