package tunnel

import (
	"context"
	"encoding/hex"
	"fmt"
	"github.com/topcloudz/fvpn/pkg/cache"
	"github.com/topcloudz/fvpn/pkg/handler"
	"github.com/topcloudz/fvpn/pkg/middleware/infra"
	"github.com/topcloudz/fvpn/pkg/option"
	"github.com/topcloudz/fvpn/pkg/packet"
	"github.com/topcloudz/fvpn/pkg/packet/handshake"
	"github.com/topcloudz/fvpn/pkg/packet/header"
	"github.com/topcloudz/fvpn/pkg/packet/notify"
	notifyack "github.com/topcloudz/fvpn/pkg/packet/notify/ack"
	peerack "github.com/topcloudz/fvpn/pkg/packet/peer/ack"
	"github.com/topcloudz/fvpn/pkg/packet/register/ack"
	"github.com/topcloudz/fvpn/pkg/socket"
	"github.com/topcloudz/fvpn/pkg/util"
	"net"
	"time"
)

// Handle union udp handler
func (t *Tunnel) Handle() handler.HandlerFunc {

	return func(ctx context.Context, frame *packet.Frame) error {
		//dest := ctx.Value("destAddr").(string)
		buff := frame.Buff[:]

		headerBuff, err := header.Decode(buff)
		if err != nil {
			return err
		}

		frame.NetworkId = hex.EncodeToString(headerBuff.NetworkId[:])
		c := ctx.Value("cache").(*cache.Cache)

		//frame.FrameType = headerBuff.Flags
		switch headerBuff.Flags {
		case option.MsgTypeRegisterAck:
			regAck, err := ack.Decode(buff)
			if err != nil {
				return err
			}
			logger.Debugf("register success, got server server ack: (%v)", regAck)
			break
		case option.MsgTypeQueryPeer:
			logger.Debugf("start get query response")
			peerPacketAck, err := peerack.Decode(buff)
			if err != nil {
				return err
			}
			infos := peerPacketAck.NodeInfos
			logger.Infof("got server peers: (%v)", infos)

			for _, info := range infos {
				//logger.Debugf("got remote node: mac: %v, ip: %s,  natIP: %s, natPort: %d", info.Mac, info.IP, info.NatIp, info.NatPort)
				address, err := util.GetAddress(info.NatIp.String(), int(info.NatPort))
				if err != nil {
					return err
				}
				node, err := c.GetNodeInfo(frame.NetworkId, info.IP.String())
				if node == nil || err != nil {
					nodeInfo := &cache.Endpoint{
						MacAddr: info.Mac,
						IP:      info.IP,
						Port:    info.Port,
						P2P:     false,
						Addr:    &address,
					}

					//cache.Nodes[info.Mac.String()] = nodeInfo
					c.SetCache(frame.NetworkId, info.IP.String(), nodeInfo)
				}
			}
			break
		case option.MsgTypePacket:
			frame.Packet = buff[:]
			t.Inbound <- frame
		case option.MsgTypeNotify:
			np, err := notify.Decode(buff)
			if err != nil {
				return err
			}
			//write back a notify
			logger.Debugf(">>>>>>>>>>>>>>got p2p notify, will create p2p tunnel........, source ip:%v, source port: %v, remote addr: %v, remote nat port: %v", np.SourceIP, np.Port, np.NatIP, np.NatPort)
			buff, err := t.buildNotifyMessageAck(np.SourceIP.String(), frame.NetworkId)
			if err != nil {
				return err
			}

			t.socket.Write(buff)
			// write back notify
			go func() {
				t.handshaking(frame, np.NatIP, int(np.NatPort), np.SourceIP.String())
			}()
		case option.MsgTypeNotifyAck:
			nck, err := notifyack.Decode(buff)
			if err != nil {
				return err
			}
			logger.Debugf("got p2p notify ack, will create p2p tunnel........, source ip:%v, source port: %v, remote addr: %v, remote nat port: %v", nck.SourceIP, nck.Port, nck.NatIP, nck.NatPort)
			go func() {
				t.handshaking(frame, nck.NatIP, int(nck.NatPort), nck.SourceIP.String())
			}()
		}

		return nil
	}
}

func (t *Tunnel) handshaking(frame *packet.Frame, natIP net.IP, natPort int, destIP string) {
	//begin to punch hole
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()
	portPair := t.manager.GetNotifyPortPair(destIP)
	logger.Debugf("===========got cached port pair, source ip: %v, source port: %v, nat ip: %v, nat port: %v", portPair.SrcIP, portPair.SrcPort, portPair.NatIP, portPair.NatPort)
	conn, err := socket.NewSocket(fmt.Sprintf("%s:%d", net.IPv4zero.String(), portPair.SrcPort), fmt.Sprintf("%s:%d", natIP.String(), natPort))
	if err != nil {
		logger.Errorf("%v", err)
		return
	}

	a := conn.LocalAddr()
	logger.Debugf(">>>>>>>>>>>>connection : %v", a)
	stopCh := make(chan int, 1)
	go func() {
		handPkt := handshake.NewPacket(frame.NetworkId)
		buff, err := handshake.Encode(handPkt)
		if err != nil {
			logger.Errorf("invalid handshake packet")
		}

		for {
			select {
			case <-stopCh:
				logger.Debugf("=============================================")
				return
			default:
				logger.Debugf("senging data to punch hole")
				_, err := conn.Write(buff)
				if err != nil {
					logger.Errorf("bad handshake packet: %v", err)
				}
				time.Sleep(time.Second * 5)
			}

		}
	}()

	go func() {
		for {
			buff := make([]byte, 1024)
			_, err = conn.Read(buff)
			if err != nil {
				logger.Errorf("punch hole failed. try again: %v", err)
				continue
			}

			//success
			logger.Debugf("punch hole success. will create a new tunnel")
			p2pTunnel := NewTunnel(t.tunHandler, conn, t.devices, infra.Middlewares(true, true), t.manager)
			t.manager.SetTunnel(destIP, p2pTunnel)
			p2pTunnel.Start() //start this p2p tunnel to service data
			stopCh <- 1
			break
		}

	}()

	select {
	case <-ctx.Done():
		logger.Debugf("punch hole finished.")
	case <-time.After(time.Second * 30):
		fmt.Println("timeout!!!")
	}
}
