package edge

import (
	"github.com/interstellar-cloud/star/pkg/option"
	"github.com/interstellar-cloud/star/pkg/packet/peer"
	"net"
)

func (edge *StarEdge) queryPeer(conn net.Conn) error {
	cp := peer.NewPacket()
	data, err := peer.Encode(cp)
	if err != nil {
		return err
	}

	switch edge.Protocol {
	case option.UDP:
		if _, err := conn.(*net.UDPConn).Write(data); err != nil {
			return nil
		}
		break
	}
	return nil
}
