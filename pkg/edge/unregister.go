package edge

import (
	"fmt"
	"github.com/interstellar-cloud/star/pkg/log"
	"github.com/interstellar-cloud/star/pkg/option"
	"github.com/interstellar-cloud/star/pkg/packet/register"
	"net"
)

// register register a edgestar to center.
func (edge *StarEdge) unregister(conn net.Conn) error {
	var err error

	rp := register.NewUnregisterPacket()
	hw, _ := net.ParseMAC(edge.MacAddr)
	rp.SrcMac = hw
	data, err := register.Encode(rp)
	fmt.Println("sending unregister data: ", data)
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
