package edge

import (
	"github.com/interstellar-cloud/star/pkg/log"
	"github.com/interstellar-cloud/star/pkg/option"
	"github.com/interstellar-cloud/star/pkg/packet/register"
	"net"
)

// register register a edgestar to center.
func (edge *StarEdge) register(conn net.Conn) error {
	var err error
	rp := register.NewPacket()
	hw, _ := net.ParseMAC(GetLocalMacAddr())
	rp.SrcMac = hw
	data, err := register.Encode(rp)
	log.Logger.Infof("sending registry data: %v", data)
	if err != nil {
		return err
	}

	switch edge.Protocol {
	case option.UDP:
		log.Logger.Infof("star start to registry self to registry: %v", rp)
		if _, err := conn.(*net.UDPConn).Write(data); err != nil {
			return err
		}
		break
	}
	return nil
}
