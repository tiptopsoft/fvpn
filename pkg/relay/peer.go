package relay

import (
	"github.com/topcloudz/fvpn/pkg/cache"
	"github.com/topcloudz/fvpn/pkg/packet/peer/ack"
)

func getPeerInfo(nodes []*cache.Endpoint) ([]ack.EdgeInfo, uint8, error) {
	var result []ack.EdgeInfo

	for _, peer := range nodes {

		info := ack.EdgeInfo{
			Mac:     peer.MacAddr,
			IP:      peer.IP,
			Port:    peer.Port,
			NatIp:   peer.NatIP,
			NatPort: peer.NatPort,
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
