package executor

import (
	"fmt"
	"github.com/interstellar-cloud/star/pkg/device"
	"github.com/interstellar-cloud/star/pkg/log"
	"github.com/interstellar-cloud/star/pkg/option"
	"github.com/interstellar-cloud/star/pkg/packet/common"
	peerack "github.com/interstellar-cloud/star/pkg/packet/peer/ack"
	"github.com/interstellar-cloud/star/pkg/packet/register/ack"
	"github.com/interstellar-cloud/star/pkg/socket"
	"io"
	"net"
)

type EdgeExecutor struct {
	Tap      *device.Tuntap
	Protocol option.Protocol
	*socket.EventLoop
}

func (ee EdgeExecutor) Execute(socket socket.Socket) error {

	if ee.Protocol == option.UDP {

		for {
			udpBytes := make([]byte, 2048)
			//_, _, err := conn.(*net.UDPConn).ReadFromUDP(udpBytes)
			_, err := socket.Read(udpBytes)
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
				log.Logger.Errorf("decode err: %v", err)
			}

			// Socket
			if socket.FileDescriptor == ee.SocketFileDescriptor {

			} else { //TAP

			}

			switch cp.Flags {
			case option.MsgTypeRegisterAck:
				regAck, err := ack.Decode(udpBytes)
				if err != nil {
					return err
				}
				log.Logger.Infof("got registry registry ack: %v", regAck)
				//create tap device
				if tap, err := device.New(device.TAP); err != nil {
					return err
				} else {
					ee.Tap = tap
					//设置IP
					if err = option.ExecCommand("/bin/sh", "-c", fmt.Sprintf("ifconfig %s %s netmask %s mtu %d up", tap.Name, regAck.AutoIP.String(), regAck.Mask.String(), 1420)); err != nil {
						return err
					}

					if err := ee.TapFd(int(tap.Fd)); err != nil {
						return err
					}
				}
				break
			case option.MsgTypePeerInfo:
				//get peerInfo
				peerPacketAck, err := peerack.Decode(udpBytes)
				if err != nil {
					return err
				}

				infos := peerPacketAck.PeerInfos
				for _, v := range infos {
					addr, err := net.ResolveUDPAddr("udp", fmt.Sprintf("%s:%d", v.Host.String(), v.Port))
					if err != nil {
						log.Logger.Errorf("resolve addr failed. err: %v", err)
					}
					option.AddrMap.Store(v.Mac.String(), addr)
				}

				break
			}

		}

	}
	return nil
}
