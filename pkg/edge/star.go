package edge

import (
	"fmt"
	"github.com/interstellar-cloud/star/pkg/device"
	"github.com/interstellar-cloud/star/pkg/option"
	"github.com/interstellar-cloud/star/pkg/pack/common"
	"github.com/interstellar-cloud/star/pkg/pack/register"
	"github.com/interstellar-cloud/star/pkg/pack/register/ack"
	"net"
)

type EdgeStar struct {
	*EdgeConfig
}

/**
 * Start logic: start to:
1. PING to super node 2. register to super 3. auto ip config tuntap 4.
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
	case TCP:
		listener, err := net.Listen("tcp", address)
		if err != nil {
			return nil, err
		}
		conn, err = listener.Accept()
	case UDP:
		conn, err = net.ListenUDP("udp", &net.UDPAddr{IP: net.IPv4(0, 0, 0, 0),
			Port: int(common.DefaultPort)})
	}

	//defer conn.Close()
	if err != nil {
		panic(err)
	}

	return conn, nil
}

// register register a edgestar to center.
func (es *EdgeStar) register(conn net.Conn) error {
	p := common.NewPacket()
	p.Flags = option.MSG_TYPE_REGISTER
	p.TTL = common.DefaultTTL

	rp := register.NewPacket()

	mac, err := option.GetLocalMac(es.TapName)
	if err != nil {
		return option.ErrGetMac
	}

	copy(rp.SrcMac[:], mac[:])

	data, err := p.Encode()
	if err != nil {
		return err
	}

	switch es.Protocol {
	case UDP:
		if _, err := conn.(*net.UDPConn).Write(data); err != nil {
			return err
		}
		break
	}
	return nil
}

func (es *EdgeStar) process(conn net.Conn) error {
	if es.Protocol == UDP {
		udpBytes := make([]byte, 2048)
		_, _, err := conn.(*net.UDPConn).ReadFromUDP(udpBytes)
		if err != nil {
			fmt.Println(err)
		}

		cp := &common.CommonPacket{}
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

			//create tap device
			if tap, err := device.New(device.TAP); err != nil {
				return err
			} else {
				//设置IP
				if err = option.ExecCommand("/bin/sh", "-c", fmt.Sprintf("ip addr add %s dev %s", regAck.AutoIP, tap.Name)); err != nil {
					return err
				}
			}

			break

		}
	}
	return nil
}
