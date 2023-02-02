package registry

import (
	"errors"
	"github.com/interstellar-cloud/star/pkg/log"
	"github.com/interstellar-cloud/star/pkg/packet/common"
	"github.com/interstellar-cloud/star/pkg/packet/peer"
	"github.com/interstellar-cloud/star/pkg/packet/peer/ack"
	"net"
)

func (r *RegStar) processPeer(addr *net.UDPAddr, conn *net.UDPConn, data []byte, cp *common.CommonPacket) {
	var p peer.PeerPacket
	p, err := peer.DecodeWithCommonPacket(data, *cp)
	if err != nil {
		log.Logger.Errorf("decode peer packet failed. err: %v", err)
	}

	// get peer info
	peers, size, err := getPeerInfo(p.SrcMac.String())
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
	<-limitChan
}

func getPeerInfo(mac string) ([]ack.PeerInfo, uint8, error) {
	m1, err := net.ParseMAC(mac)
	if err != nil {
		return nil, 0, errors.New("invalid mac")
	}
	var result []ack.PeerInfo
	res, ok := m.Load(mac)
	if !ok {
		return nil, 0, errors.New("peer not found")
	}

	a := res.(*net.UDPAddr)

	info := ack.PeerInfo{
		Mac:  m1,
		Host: a.IP,
		Port: uint16(a.Port),
	}
	result = append(result, info)
	return result, 1, nil
}

func peerAckBuild(infos []ack.PeerInfo, size uint8) ([]byte, error) {
	peerPacket := ack.NewPacket()
	peerPacket.Size = size
	peerPacket.PeerInfos = infos

	return ack.Encode(peerPacket)
}
