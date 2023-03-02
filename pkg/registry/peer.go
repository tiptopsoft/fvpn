package registry

import (
	"github.com/interstellar-cloud/star/pkg/node"
	"github.com/interstellar-cloud/star/pkg/packet/peer/ack"
	"golang.org/x/sys/unix"
)

func (r *RegStar) processFindPeer(addr unix.Sockaddr) {
	logger.Infof("start to process query peers...")
	// get peer info
	peers, size, err := getPeerInfo(r.cache)
	logger.Infof("registry peers: (%v), size: (%v)", peers, size)
	if err != nil {
		logger.Errorf("get peers from registry failed. err: %v", err)
	}

	f, err := peerAckBuild(peers, size)
	if err != nil {
		logger.Errorf("get peer ack from registry failed. err: %v", err)
	}

	err = r.socket.WriteToUdp(f, addr)
	logger.Infof("addr: %v", addr)
	if err != nil {
		logger.Errorf("registry write failed. err: %v", err)
	}

	logger.Infof("finish process query peers")
}

func getPeerInfo(peers node.NodesCache) ([]ack.EdgeInfo, uint8, error) {
	var result []ack.EdgeInfo
	for _, peer := range peers.Nodes {
		info := ack.EdgeInfo{
			Mac:  peer.MacAddr,
			Host: peer.IP,
			Port: peer.Port,
		}
		result = append(result, info)
	}

	return result, uint8(len(result)), nil
}

func peerAckBuild(infos []ack.EdgeInfo, size uint8) ([]byte, error) {
	peerPacket := ack.NewPacket()
	peerPacket.Size = size
	peerPacket.PeerInfos = infos

	return peerPacket.Encode()
}
