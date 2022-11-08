package edge

import (
	"fmt"
	"github.com/interstellar-cloud/star/pkg/log"
	"github.com/interstellar-cloud/star/pkg/option"
	"github.com/interstellar-cloud/star/pkg/packet/common"
	"github.com/interstellar-cloud/star/pkg/packet/register"
	"github.com/interstellar-cloud/star/pkg/packet/register/ack"
	"net"
)

type EdgeStar struct {
	*option.EdgeConfig
}

/**
 * Start logic: start to:
1. PING to register node 2. register to register 3. auto ip config tuntap 4.
*/
func (edge EdgeStar) Start() error {
	//init connect to registry
	var conn net.Conn
	var err error
	conn, err = edge.conn(edge.Registry)
	if err != nil {
		return err
	}

	// registry to registry
	err = edge.register(conn)
	if err != nil {
		return err
	}

	//run loop process udp
	err = edge.process(conn)
	if err != nil {
		return err
	}
	return nil
}

func (es *EdgeStar) conn(address string) (net.Conn, error) {
	var conn net.Conn
	var err error

	switch es.Protocol {
	case option.UDP:
		conn, err = net.Dial("udp", es.Registry)
	}

	//defer conn.Close()
	if err != nil {
		return nil, err
	}

	log.Logger.Info("star connected to server: %s", es.Registry)
	return conn, nil
}

// register register a edgestar to center.
func (es *EdgeStar) register(conn net.Conn) error {
	var err error
	p := common.NewPacket()
	p.Flags = option.MSG_TYPE_REGISTER_SUPER
	p.TTL = common.DefaultTTL

	rp := register.NewPacket()

	mac := es.MacAddr
	//if mac == "" {
	//	mac, err = option.GetLocalMac(es.TapName)
	//	if err != nil {
	//		return option.ErrGetMac
	//	}
	//}

	copy(rp.SrcMac[:], mac[:])

	data, err := p.Encode()
	if err != nil {
		return err
	}

	switch es.Protocol {
	case option.UDP:
		log.Logger.Infof("star start to register self to registry: %v", rp)
		if _, err := conn.(*net.UDPConn).Write(data); err != nil {
			return err
		}
		break
	}
	return nil
}

func (es *EdgeStar) process(conn net.Conn) error {
	if es.Protocol == option.UDP {
		udpBytes := make([]byte, 2048)
		_, _, err := conn.(*net.UDPConn).ReadFromUDP(udpBytes)
		if err != nil {
			fmt.Println(err)
		}

		cp := common.CommonPacket{}
		cp, err = cp.Decode(udpBytes)

		if err != nil {
			fmt.Println(err)
		}

		switch cp.Flags {
		case option.MSG_TYPE_REGISTER_ACK:
			regAck, err := ack.NewPacket().Decode(udpBytes)
			if err != nil {
				return err
			}
			log.Logger.Infof("got registry register ack: %v", regAck)
			//create tap device
			//if tap, err := device.New(device.TAP); err != nil {
			//	return err
			//} else {
			//	//设置IP
			//	address := fmt.Sprintf("%d:%d:%d:%d", regAck.AutoIP[0], regAck.AutoIP[1], regAck.AutoIP[2], regAck.AutoIP[3])
			//	if err = option.ExecCommand("/bin/sh", "-c", fmt.Sprintf("ip addr add %s dev %s", address, tap.Name)); err != nil {
			//		return err
			//	}
			//}

			break

		}
	}
	return nil
}
