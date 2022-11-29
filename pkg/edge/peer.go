package edge

import (
	"github.com/interstellar-cloud/star/pkg/option"
	"github.com/interstellar-cloud/star/pkg/packet/common"
	"net"
)

func (es *EdgeStar) queryPeer(conn net.Conn) error {
	cp := common.NewPacket()
	cp.Flags = option.MsgTypePeerInfo

	data, err := common.Encode(cp)
	if err != nil {
		return err
	}

	switch es.Protocol {
	case option.UDP:
		if _, err := conn.(*net.UDPConn).Write(data); err != nil {
			return nil
		}
		break
	}
	return nil
}
