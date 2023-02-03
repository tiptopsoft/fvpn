package registry

import (
	"github.com/interstellar-cloud/star/pkg/util"
	"github.com/interstellar-cloud/star/pkg/util/log"
	"github.com/interstellar-cloud/star/pkg/util/packet/common"
	"github.com/interstellar-cloud/star/pkg/util/packet/forward"
	"net"
	"sync"
)

var socketMap sync.Map

func (r *RegStar) forward(data []byte, cp *common.CommonPacket) {

	fp, err := forward.Decode(data)
	if err != nil {
		log.Logger.Errorf("decode forward packet failed. err: %v", err)
	}

	// find Addr in registry
	peer := util.FindPeers(r.Peers, fp.DstMac.String())

	if peer == nil {
		log.Logger.Errorf("dst has not registerd in registry. macAddr: %s", fp.DstMac.String())
	} else if util.IsBroadCast(fp.DstMac.String()) {
		//broad cast send data to all edge
		for _, v := range r.Peers {
			sock := v.Conn
			_, err := sock.(*net.UDPConn).Write(data)
			if err != nil {
				log.Logger.Errorf("send to remote edge or registry failed. err: %v", err)
			}
		}

	}

}
