package registry

import (
	"github.com/interstellar-cloud/star/pkg/util"
	"github.com/interstellar-cloud/star/pkg/util/log"
	"github.com/interstellar-cloud/star/pkg/util/node"
	"github.com/interstellar-cloud/star/pkg/util/packet/common"
	"github.com/interstellar-cloud/star/pkg/util/packet/forward"
)

func (r *RegStar) forward(data []byte, cp *common.CommonPacket) {
	log.Logger.Infof("registry got forward packet: %v", data)
	fp, err := forward.Decode(data)
	if err != nil {
		log.Logger.Errorf("decode forward packet failed. err: %v", err)
	}

	//if util.IsBroadCast(fp.DstMac.String()) {
	peer := node.FindNode(r.cache, fp.DstMac.String())
	if peer == nil {
		//用dstIP去查询
		ip := util.GetDstIP(data)
		peer = node.FindNodeByIP(r.cache, ip.String())
		if peer == nil {
			log.Logger.Errorf("dst has not registerd in registry. macAddr: %s, addr: %s", fp.DstMac.String())
		}

		return
	}

	for _, v := range r.cache.Nodes {
		err := r.socket.WriteToUdp(data, v.Addr)
		log.Logger.Infof("forward packet: (%v), addr: %v", data, v.Addr)
		if err != nil {
			log.Logger.Errorf("send to remote edge or registry failed. err: %v", err)
		}
	}

	//} else {
	//	// find Addr in registry
	//	peer := util.FindPeers(r.cache, fp.DstMac.String())
	//	if peer == nil {
	//		log.Logger.Errorf("dst has not registerd in registry. macAddr: %s", fp.DstMac.String())
	//	}
	//}
}
