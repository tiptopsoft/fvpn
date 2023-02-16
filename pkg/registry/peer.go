package registry

import (
	"github.com/interstellar-cloud/star/pkg/util/log"
	"github.com/interstellar-cloud/star/pkg/util/node"
	"github.com/interstellar-cloud/star/pkg/util/packet/peer/ack"
	"golang.org/x/sys/unix"
)

func (r *RegStar) processFindPeer(addr unix.Sockaddr) {
	log.Logger.Infof("start to process query peers...")
	// get peer info
	peers, size, err := getPeerInfo(r.cache)
	log.Logger.Infof("registry peers: (%v), size: (%v)", peers, size)
	if err != nil {
		log.Logger.Errorf("get peers from registry failed. err: %v", err)
	}

	f, err := peerAckBuild(peers, size)
	if err != nil {
		log.Logger.Errorf("get peer ack from registry failed. err: %v", err)
	}

	err = r.socket.WriteToUdp(f, addr)
	log.Logger.Infof("addr: %v", addr)
	if err != nil {
		log.Logger.Errorf("registry write failed. err: %v", err)
	}

	log.Logger.Infof("finish process query peers")
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

	return ack.Encode(peerPacket)
}
