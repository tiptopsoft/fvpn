package registry

import (
	"github.com/interstellar-cloud/star/pkg/packet/peer/ack"
	"github.com/interstellar-cloud/star/pkg/util/log"
	"net"
)

func (r *RegStar) processPeer(addr *net.UDPAddr, conn *net.UDPConn) {
	log.Logger.Infof("start to process query peers...")
	// get peer info
	peers, size, err := getPeerInfo()
	log.Logger.Infof("registry peers: (%v), size: (%v)", peers, size)
	if err != nil {
		log.Logger.Errorf("get peers from registry failed. err: %v", err)
	}

	f, err := peerAckBuild(peers, size)
	if err != nil {
		log.Logger.Errorf("get peer ack from registry failed. err: %v", err)
	}

	_, err = conn.WriteToUDP(f, addr)
	if err != nil {
		log.Logger.Errorf("registry write failed. err: %v", err)
	}

	log.Logger.Infof("finish process query peers. (%v)", f)
}

func getPeerInfo() ([]ack.PeerInfo, uint8, error) {
	var result []ack.PeerInfo
	m.Range(func(mac, pub any) bool {
		a := pub.(*net.UDPAddr)
		info := ack.PeerInfo{
			Mac:  []byte(mac.(string)),
			Host: a.IP,
			Port: uint16(a.Port),
		}
		result = append(result, info)
		return true
	})

	return result, uint8(len(result)), nil
}

func peerAckBuild(infos []ack.PeerInfo, size uint8) ([]byte, error) {
	peerPacket := ack.NewPacket()
	peerPacket.Size = size
	peerPacket.PeerInfos = infos

	return ack.Encode(peerPacket)
}
