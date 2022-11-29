package edge

import (
	"fmt"
	"github.com/interstellar-cloud/star/pkg/log"
	"github.com/interstellar-cloud/star/pkg/option"
	"github.com/interstellar-cloud/star/pkg/packet/common"
	"github.com/interstellar-cloud/star/pkg/packet/register"
	"net"
)

// register register a edgestar to center.
func (es *EdgeStar) unregister(conn net.Conn) error {
	var err error
	p := common.NewPacket()
	p.Flags = option.MsgTypeUnregisterSuper
	p.TTL = common.DefaultTTL
	rp := register.NewPacket()
	hw, _ := net.ParseMAC(es.MacAddr)

	rp.SrcMac = hw
	rp.CommonPacket = p
	data, err := register.Encode(rp)
	fmt.Println("sending unregister data: ", data)
	if err != nil {
		return err
	}

	switch es.Protocol {
	case option.UDP:
		log.Logger.Infof("star start to registry self to registry: %v", rp)
		if _, err := conn.(*net.UDPConn).Write(data); err != nil {
			return err
		}
		break
	}
	return nil
}
