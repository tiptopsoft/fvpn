package edge

import (
	"github.com/interstellar-cloud/star/pkg/log"
	"github.com/interstellar-cloud/star/pkg/option"
	"github.com/interstellar-cloud/star/pkg/packet/common"
	"github.com/interstellar-cloud/star/pkg/packet/register"
	"net"
)

// register register a edgestar to center.
func (es *EdgeStar) register(conn net.Conn) error {
	var err error
	p := common.NewPacket()
	p.Flags = option.MsgTypeRegisterSuper
	p.TTL = common.DefaultTTL
	rp := register.NewPacket()
	hw, _ := net.ParseMAC(GetLocalMacAddr())

	rp.SrcMac = hw
	rp.CommonPacket = p
	data, err := register.Encode(rp)
	log.Logger.Infof("sending data: %v", data)
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
