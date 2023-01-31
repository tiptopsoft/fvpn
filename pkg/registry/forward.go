package registry

import (
	"github.com/interstellar-cloud/star/pkg/log"
	"github.com/interstellar-cloud/star/pkg/packet/common"
	"github.com/interstellar-cloud/star/pkg/packet/forward"
	"net"
	"sync"
)

var once sync.Once
var socketMap sync.Map

func (r *RegStar) forward(data []byte, cp *common.CommonPacket) {

	fp, err := forward.Decode(data)
	if err != nil {
		log.Logger.Errorf("decode forward packet failed. err: %v", err)
	}

	// find Addr in registry

	if addr, ok := m.Load(fp.DstMac); !ok {
		log.Logger.Errorf("dst has not registerd in registry. addr: %s", addr)
	} else {
		if sock, ok := socketMap.Load(fp.DstMac); !ok {
			once.Do(func() {

				conn, err := net.Dial("udp", addr.(*net.UDPAddr).String())
				if err != nil {
					log.Logger.Errorf("dial remote edge failed. err: %v", err)
				}
				sock = conn
				socketMap.Store(fp.DstMac, conn)
			})
		} else {
			_, err := sock.(*net.UDPConn).Write(data)
			if err != nil {
				log.Logger.Errorf("send to remote edge failed. err: %v", err)
			}
		}

	}

}
