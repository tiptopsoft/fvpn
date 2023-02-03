package registry

import (
	"encoding/json"
	"github.com/interstellar-cloud/star/pkg/util"
	"github.com/interstellar-cloud/star/pkg/util/log"
	"github.com/interstellar-cloud/star/pkg/util/packet/peer/ack"
	"github.com/interstellar-cloud/star/pkg/util/socket"
	"net"
)

func (r *RegStar) processFindPeer(addr *net.UDPAddr, socket socket.Socket) {
	log.Logger.Infof("start to process query peers...")
	// get peer info
	peers, size, err := getPeerInfo(r.Peers)
	b, err := json.Marshal(peers)
	if err != nil {
		log.Logger.Errorf("proces peer failed: %v", err)
	}
	log.Logger.Infof("registry peers: (%v), size: (%v)", string(b), size)
	if err != nil {
		log.Logger.Errorf("get peers from registry failed. err: %v", err)
	}

	f, err := peerAckBuild(peers, size)
	if err != nil {
		log.Logger.Errorf("get peer ack from registry failed. err: %v", err)
	}

	_, err = socket.WriteToUdp(f, addr)
	log.Logger.Infof("addr: %v", addr)
	if err != nil {
		log.Logger.Errorf("registry write failed. err: %v", err)
	}

	log.Logger.Infof("finish process query peers. (%v)", f)
}

func getPeerInfo(peers util.Peers) ([]ack.EdgeInfo, uint8, error) {
	var result []ack.EdgeInfo
	for _, peer := range peers {
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
