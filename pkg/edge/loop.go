package edge

import (
	"fmt"
	"github.com/interstellar-cloud/star/pkg/device"
	"github.com/interstellar-cloud/star/pkg/log"
	"github.com/interstellar-cloud/star/pkg/option"
	"github.com/interstellar-cloud/star/pkg/packet/common"
	"github.com/interstellar-cloud/star/pkg/packet/register/ack"
	"io"
	"net"
)

func (es *EdgeStar) process(conn net.Conn) error {
	if es.Protocol == option.UDP {
		for {
			udpBytes := make([]byte, 2048)
			_, _, err := conn.(*net.UDPConn).ReadFromUDP(udpBytes)
			if err != nil {
				if err == io.EOF {
					//no data exists, continue read next frame.
					continue
				} else {
					log.Logger.Errorf("read from remote error: %v", err)
				}
			}

			cp, err := common.Decode(udpBytes)

			if err != nil {
				fmt.Println(err)
			}

			switch cp.Flags {
			case option.MSG_TYPE_REGISTER_ACK:
				regAck, err := ack.Decode(udpBytes)
				if err != nil {
					return err
				}
				log.Logger.Infof("got registry register ack: %v", regAck)
				//create tap device
				if tap, err := device.New(device.TAP); err != nil {
					return err
				} else {
					es.tap = tap
					//设置IP
					address := fmt.Sprintf("%d:%d:%d:%d", regAck.AutoIP[0], regAck.AutoIP[1], regAck.AutoIP[2], regAck.AutoIP[3])
					if err = option.ExecCommand("/bin/sh", "-c", fmt.Sprintf("ip addr add %s dev %s", address, tap.Name)); err != nil {
						return err
					}
				}
				ch <- 2
				break
			case option.MSG_TYPE_PEER_INFO:
				break
			}

		}

	}
	return nil
}