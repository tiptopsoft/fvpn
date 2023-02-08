package registry

import (
	"github.com/interstellar-cloud/star/pkg/util/log"
	"github.com/interstellar-cloud/star/pkg/util/packet/common"
	"sync"
)

var socketMap sync.Map

func (r *RegStar) forward(data []byte, cp *common.CommonPacket) {
	log.Logger.Infof("registry got forward packet: %v", data)
	//fp, err := forward.Decode(data)
	//if err != nil {
	//	log.Logger.Errorf("decode forward packet failed. err: %v", err)
	//}

	//if util.IsBroadCast(fp.DstMac.String()) {
	//broad cast send data to all edge
	for _, v := range r.Nodes {
		err := r.socket.WriteToUdp(data, v.Addr)
		log.Logger.Infof("forward packet: (%v), addr: %v", data, v.Addr)
		if err != nil {
			log.Logger.Errorf("send to remote edge or registry failed. err: %v", err)
		}
	}

	//} else {
	//	// find Addr in registry
	//	peer := util.FindPeers(r.Nodes, fp.DstMac.String())
	//	if peer == nil {
	//		log.Logger.Errorf("dst has not registerd in registry. macAddr: %s", fp.DstMac.String())
	//	}
	//}
}
