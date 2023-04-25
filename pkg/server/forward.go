package server

import (
	"github.com/topcloudz/fvpn/pkg/cache"
	"github.com/topcloudz/fvpn/pkg/packet"
	"github.com/topcloudz/fvpn/pkg/packet/forward"
	"github.com/topcloudz/fvpn/pkg/util"
)

func (r *RegStar) forward(data []byte, cp *packet.Header) {
	logger.Infof("server got forward packet: %v", data)
	fpInterface, err := forward.NewPacket("").Decode(data)
	fp := fpInterface.(forward.ForwardPacket)

	if err != nil {
		logger.Errorf("decode forward packet failed. err: %v", err)
	}

	//if util.IsBroadCast(fp.DstMac.String()) {
	peer := cache.FindPeer(r.cache, fp.DstMac.String())
	if peer == nil {
		//用dstIP去查询
		ip := util.GetDstIP(data)
		peer = cache.FindPeerByIP(r.cache, ip.String())
		if peer == nil {
			logger.Errorf("dst has not registerd in server. macAddr: %s, addr: %s", fp.DstMac.String())
		}

		return
	}

	for _, v := range r.cache.Nodes {
		err := r.socket.WriteToUdp(data, v.Addr)
		logger.Infof("forward packet: (%v), addr: %v", data, v.Addr)
		if err != nil {
			logger.Errorf("send to remote client or server failed. err: %v", err)
		}
	}

	//} else {
	//	// find Addr in server
	//	peer := util.FindPeers(r.cache, fp.DstMac.String())
	//	if peer == nil {
	//		logger.Errorf("dst has not registerd in server. macAddr: %s", fp.DstMac.String())
	//	}
	//}
}
