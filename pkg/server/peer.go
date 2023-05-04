package server

import (
	"github.com/topcloudz/fvpn/pkg/cache"
	"github.com/topcloudz/fvpn/pkg/packet/peer/ack"
)

func getPeerInfo(nodes []*cache.NodeInfo) ([]ack.EdgeInfo, uint8, error) {
	var result []ack.EdgeInfo

	for _, peer := range nodes {
		info := ack.EdgeInfo{
			Mac:  peer.MacAddr,
			IP:   peer.IP,
			Port: peer.Port,
		}
		result = append(result, info)
	}

	return result, uint8(len(result)), nil
}

func peerAckBuild(infos []ack.EdgeInfo, size uint8) ([]byte, error) {
	peerPacket := ack.NewPacket()
	peerPacket.Size = size
	peerPacket.NodeInfos = infos

	return ack.Encode(peerPacket)
}
