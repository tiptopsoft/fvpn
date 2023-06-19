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
	peerack "github.com/topcloudz/fvpn/pkg/packet/peer/ack"
	"github.com/topcloudz/fvpn/pkg/packet/register/ack"
	"github.com/topcloudz/fvpn/pkg/socket"
	"github.com/topcloudz/fvpn/pkg/util"
	"golang.org/x/sys/unix"
	"time"
)

// Handle union udp handler
func (t *Tunnel) Handle() handler.HandlerFunc {

	return func(ctx context.Context, frame *packet.Frame) error {
		dest := ctx.Value("destAddr").(string)
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
					sock := socket.NewSocket(6061)
					nodeInfo := &cache.Endpoint{
						Socket:  sock,
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
		case option.MsgTypeNotifyType:
			np, err := notify.Decode(buff)
			if err != nil {
				return err
			}

			//use socket write a notify to dest
			buff, err := buildNotifyMessage(frame.NetworkId)
			if err != nil {
				return err
			}

			// write a handshake
			t.socket.Write(buff)

			//begin to punch hole
			ctx, cancel := context.WithTimeout(context.Background(), time.Duration(time.Second*30))
			defer cancel()
			go func() {
				portPair := <-Pool.ch
				conn := socket.NewSocket(int(portPair.SrcPort))
				destAddr := unix.SockaddrInet4{Port: int(np.NatPort)}
				copy(destAddr.Addr[:], portPair.NatIP.To4())
				conn.Connect(&destAddr)

				handPkt := handshake.NewPacket(frame.NetworkId)
				buff, err := handshake.Encode(handPkt)
				if err != nil {
					logger.Errorf("bad handshake packet")
					return
				}

				for {
					_, err := conn.Write(buff)
					if err != nil {
						logger.Errorf("bad handshake packet")
						return
					}

					buff := make([]byte, 1024)
					_, err = conn.Read(buff)
					if err != nil {
						logger.Errorf("punch hole failed. try again")
						continue
					}

					//success
					logger.Debugf("punch hole success. will create a new tunnel")
					p2pTunnel := NewTunnel(t.tunHandler, conn, t.devices, infra.Middlewares(true, true), t.manager)
					t.manager.SetTunnel(dest, p2pTunnel)
					p2pTunnel.Start() //start this p2p tunnel to service data
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

		return nil
	}

}
