package server

import (
	"github.com/topcloudz/fvpn/pkg/cache"
	"github.com/topcloudz/fvpn/pkg/packet/peer/ack"
	"golang.org/x/sys/unix"
)

func (r *RegStar) processFindPeer(addr unix.Sockaddr) {
	logger.Infof("start to process query peers...")
	// get peer info
	peers, size, err := getPeerInfo(r.cache.GetNodes())
	logger.Infof("server peers: (%v), size: (%v)", peers, size)
	if err != nil {
		logger.Errorf("get peers from server failed. err: %v", err)
	}

	f, err := peerAckBuild(peers, size)
	if err != nil {
		logger.Errorf("get peer ack from server failed. err: %v", err)
	}

	err = r.socket.WriteToUdp(f, addr)
	logger.Infof("addr: %v", addr)
	if err != nil {
		logger.Errorf("server write failed. err: %v", err)
	}

	logger.Infof("finish process query peers")
}

func getPeerInfo(nodes []*cache.NodeInfo) ([]ack.EdgeInfo, uint8, error) {
	var result []ack.EdgeInfo

	for _, peer := range nodes {
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
