package edge

import (
	"errors"
	"fmt"
	"github.com/interstellar-cloud/star/pkg/util"
	"github.com/interstellar-cloud/star/pkg/util/addr"
	"github.com/interstellar-cloud/star/pkg/util/log"
	"github.com/interstellar-cloud/star/pkg/util/option"
	"github.com/interstellar-cloud/star/pkg/util/packet"
	"github.com/interstellar-cloud/star/pkg/util/packet/common"
	"github.com/interstellar-cloud/star/pkg/util/packet/forward"
	peerack "github.com/interstellar-cloud/star/pkg/util/packet/peer/ack"
	"github.com/interstellar-cloud/star/pkg/util/packet/register/ack"
	"github.com/interstellar-cloud/star/pkg/util/socket"
	"github.com/interstellar-cloud/star/pkg/util/tuntap"
	"io"
	"net"
	"sync"
	"unsafe"
)

var (
	once sync.Once
)

type Star struct {
	*option.StarConfig
	tuntap *tuntap.Tuntap
	socket.Socket
	Nodes      util.Nodes //获取回来的Peers  mac: Node
	socketFunc func(device *tuntap.Tuntap, socket socket.Socket) error
	tapFunc    func(device *tuntap.Tuntap, socket socket.Socket) error
	inbound    []chan *packet.Packet
}

func (star Star) Start() error {
	once.Do(func() {
		star.Socket = socket.NewSocket()
		if err := star.conn(); err != nil {

		}
		star.Nodes = make(util.Nodes, 1)
		star.Protocol = option.UDP
		tap, err := tuntap.New(tuntap.TAP)
		star.tuntap = tap

		star.tapFunc = func(device *tuntap.Tuntap, skt socket.Socket) error {
			b := make([]byte, option.STAR_PKT_BUFF_SIZE)
			n, err := device.Read(b)
			destMac := util.GetMacAddr(b)
			fmt.Println(fmt.Sprintf("Read %d bytes from device %s, will write to dest %s", n, star.tuntap.Name, destMac))
			if err != nil {
				log.Logger.Errorf("tap read failed. (%v)", err)
				return err
			}
			broad := addr.IsBroadCast(destMac)
			//broad frame, go through supernode
			fp := forward.NewPacket()
			fp.SrcMac, err = addr.GetMacAddrByDev(star.tuntap.Name)
			if err != nil {
				log.Logger.Errorf("get src mac failed, err: %v", err)
			}
			fp.DstMac, err = net.ParseMAC(destMac)
			if err != nil {
				log.Logger.Errorf("get src mac failed, err: %v", err)
			}

			bs, err := forward.Encode(fp)
			if err != nil {
				log.Logger.Errorf("encode forward failed, err: %v", err)
			}

			idx := 0
			newPacket := make([]byte, 2048)
			idx = packet.EncodeBytes(newPacket, bs, idx)
			packet.EncodeBytes(newPacket, b[:n], idx)
			if broad {
				write2Net(skt, newPacket)
			} else {
				// go p2p
				log.Logger.Infof("find peer in edge, destMac: %v", destMac)
				p := util.FindNode(star.Nodes, destMac)
				if p == nil {
					return errors.New("peer not found, may be not registered in registry")
				}
				write2Net(p.Socket, newPacket)
			}
			return nil
		}

		star.socketFunc = func(device *tuntap.Tuntap, skt socket.Socket) error {
			fmt.Println("exec net skt.")
			if star.Protocol == option.UDP {
				udpBytes := make([]byte, 2048)
				n, err := skt.Read(udpBytes)
				log.Logger.Infof("star net skt receive size: %d, data: (%v)", n, udpBytes)
				if err != nil {
					if err == io.EOF {
						//no data exists, continue read next frame continue
						log.Logger.Errorf("not data exists")
					} else {
						log.Logger.Errorf("read from remote error: %v", err)
					}
				}

				cp, err := common.Decode(udpBytes)
				if err != nil {
					log.Logger.Errorf("decode err: %v", err)
				}

				switch cp.Flags {
				case option.MsgTypeRegisterAck:
					regAck, err := ack.Decode(udpBytes)
					if err != nil {
						return err
					}
					log.Logger.Infof("got registry registry ack: (%v)", regAck)
					//设置IP
					if err = option.ExecCommand("/bin/sh", "-c", fmt.Sprintf("ifconfig %s %s netmask %s mtu %d up", star.tuntap.Name, regAck.AutoIP.String(), regAck.Mask.String(), 1420)); err != nil {
						return err
					}
					break
				case option.MsgTypeQueryPeer:
					peerPacketAck, err := peerack.Decode(udpBytes)
					if err != nil {
						return err
					}
					infos := peerPacketAck.PeerInfos
					log.Logger.Infof("got registry peers: (%v)", infos)
					for _, info := range infos {
						address, err := util.GetAddress(info.Host.String(), int(info.Port))
						if err != nil {
							log.Logger.Errorf("resolve addr failed, err: %v", err)
						}
						sock := socket.NewSocket()
						err = sock.Connect(&address)
						if err != nil {
							return err
						}
						peerInfo := &util.Node{
							Socket:  sock,
							MacAddr: info.Mac,
							IP:      info.Host,
							Port:    info.Port,
						}
						star.Nodes[info.Mac.String()] = peerInfo
					}
					break
				case option.MsgTypePacket:
					forwardPacket, err := forward.Decode(udpBytes)
					if err != nil {
						return err
					}
					log.Logger.Infof("got through packet: %v, srcMac: %v, current tap macAddr: %v", forwardPacket, forwardPacket.SrcMac, star.tuntap.MacAddr)

					if forwardPacket.SrcMac.String() == star.tuntap.MacAddr.String() {
						//self, drop packet
						log.Logger.Warnf("self packet droped: %v, srcMac: %v, current tap macAddr: %v", forwardPacket, forwardPacket.SrcMac, star.tuntap.MacAddr)
					} else {
						//写入到tap
						idx := unsafe.Sizeof(forwardPacket)
						if _, err := star.tuntap.Write(udpBytes[idx:n]); err != nil {
							log.Logger.Errorf("write to tap failed. (%v)", err.Error())
						}
						log.Logger.Infof("net write to tap as tap response to client. size: %d", n-int(idx))
					}
					break
				}
			}
			return nil
		}

		if err != nil {
			log.Logger.Errorf("create or connect tuntap failed, err: (%v)", err)
		}

		if err := star.register(); err != nil {
			log.Logger.Errorf("registry failed. (%v)", err)
		}
	})
	star.starLoop()
	return nil
}
