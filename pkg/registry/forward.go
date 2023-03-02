package registry

import (
	"github.com/interstellar-cloud/star/pkg/node"
	"github.com/interstellar-cloud/star/pkg/packet/common"
	"github.com/interstellar-cloud/star/pkg/packet/forward"
	"github.com/interstellar-cloud/star/pkg/util"
)

func (r *RegStar) forward(data []byte, cp *common.PacketHeader) {
	logger.Infof("registry got forward packet: %v", data)
	fpInterface, err := forward.NewPacket().Decode(data)
	fp := fpInterface.(forward.ForwardPacket)

	if err != nil {
		logger.Errorf("decode forward packet failed. err: %v", err)
	}

	//if util.IsBroadCast(fp.DstMac.String()) {
	peer := node.FindNode(r.cache, fp.DstMac.String())
	if peer == nil {
		//用dstIP去查询
		ip := util.GetDstIP(data)
		peer = node.FindNodeByIP(r.cache, ip.String())
		if peer == nil {
			logger.Errorf("dst has not registerd in registry. macAddr: %s, addr: %s", fp.DstMac.String())
		}

		return
	}

	for _, v := range r.cache.Nodes {
		err := r.socket.WriteToUdp(data, v.Addr)
		logger.Infof("forward packet: (%v), addr: %v", data, v.Addr)
		if err != nil {
			logger.Errorf("send to remote edge or registry failed. err: %v", err)
		}
	}

	//} else {
	//	// find Addr in registry
	//	peer := util.FindPeers(r.cache, fp.DstMac.String())
	//	if peer == nil {
	//		logger.Errorf("dst has not registerd in registry. macAddr: %s", fp.DstMac.String())
	//	}
	//}
}
