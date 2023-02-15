package edge

import (
	"errors"
	"fmt"
	"github.com/interstellar-cloud/star/pkg/util"
	"github.com/interstellar-cloud/star/pkg/util/addr"
	"github.com/interstellar-cloud/star/pkg/util/log"
	"github.com/interstellar-cloud/star/pkg/util/node"
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
	tap *tuntap.Tuntap
	socket.Socket
	cache      node.NodesCache //获取回来的Peers  mac: Node
	socketFunc func(device *tuntap.Tuntap, socket socket.Socket) error
	tapFunc    func(device *tuntap.Tuntap, socket socket.Socket) error
	inbound    []chan *packet.Packet
}

func (star Star) Start() error {
	once.Do(func() {
		star.Socket = socket.NewSocket()
		if err := star.conn(); err != nil {

		}
		star.cache = node.New()
		star.Protocol = option.UDP
		tap, err := tuntap.New(tuntap.TAP)
		star.tap = tap

		star.tapFunc = func(device *tuntap.Tuntap, skt socket.Socket) error {
			b := make([]byte, option.STAR_PKT_BUFF_SIZE)
			size, err := device.Read(b)
			destMac := util.GetMacAddr(b)
			fmt.Println(fmt.Sprintf("Read %d bytes from device %s, will write to dest %s", size, star.tap.Name, destMac))
			if err != nil {
				log.Logger.Errorf("tap read failed. (%v)", err)
				return err
			}
			broad := addr.IsBroadCast(destMac)
			//broad frame, go through supernode
			fp := forward.NewPacket()
			fp.SrcMac, err = addr.GetMacAddrByDev(star.tap.Name)
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
			idx = packet.EncodeBytes(newPacket, b[:size], idx)
			if broad {
				write2Net(skt, newPacket[:idx])
			} else {
				// go p2p
				log.Logger.Infof("find peer in edge, destMac: %v", destMac)
				p := node.FindNode(star.cache, destMac)
				if p == nil {
					write2Net(skt, newPacket[:idx])
					log.Logger.Warnf("peer not found, go through super node")
				} else {
					write2Net(p.Socket, newPacket[:idx])
				}
			}
			return nil
		}

		star.socketFunc = func(device *tuntap.Tuntap, skt socket.Socket) error {
			if star.Protocol == option.UDP {
				udpBytes := make([]byte, 2048)
				size, err := skt.Read(udpBytes)
				if size < 0 {
					return errors.New("no data exists")
				}
				log.Logger.Infof("star net skt receive size: %d, data: (%v)", size, udpBytes[:size])
				if err != nil {
					if err == io.EOF {
						//no data exists, continue read next frame continue
						log.Logger.Errorf("not data exists")
					} else {
						log.Logger.Errorf("read from remote error: %v", err)
					}
				}

				cp, err := common.Decode(udpBytes[:size])
				if err != nil {
					log.Logger.Errorf("decode err: %v", err)
				}

				switch cp.Flags {
				case option.MsgTypeRegisterAck:
					regAck, err := ack.Decode(udpBytes[:size])
					if err != nil {
						return err
					}
					log.Logger.Infof("got registry registry ack: (%v)", regAck)
					//设置IP
					if err = option.ExecCommand("/bin/sh", "-c", fmt.Sprintf("ifconfig %s %s netmask %s mtu %d up", star.tap.Name, regAck.AutoIP.String(), regAck.Mask.String(), 1420)); err != nil {
						return err
					}
					break
				case option.MsgTypeQueryPeer:
					peerPacketAck, err := peerack.Decode(udpBytes[:size])
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
						peerInfo := &node.Node{
							Socket:  sock,
							MacAddr: info.Mac,
							IP:      info.Host,
							Port:    info.Port,
						}
						star.cache.Nodes[info.Mac.String()] = peerInfo
					}
					break
				case option.MsgTypePacket:
					forwardPacket, err := forward.Decode(udpBytes[:size])
					if err != nil {
						return err
					}
					log.Logger.Infof("got through packet: %v, srcMac: %v, current tap macAddr: %v", forwardPacket, forwardPacket.SrcMac, star.tap.MacAddr)

					if forwardPacket.SrcMac.String() == star.tap.MacAddr.String() {
						//self, drop packet
						log.Logger.Infof("self packet droped: %v, srcMac: %v, current tap macAddr: %v", forwardPacket, forwardPacket.SrcMac, star.tap.MacAddr)
					} else {
						//写入到tap
						idx := unsafe.Sizeof(forwardPacket)
						if _, err := star.tap.Write(udpBytes[idx:size]); err != nil {
							log.Logger.Errorf("write to tap failed. (%v)", err.Error())
						}
						log.Logger.Infof("net write to tap as tap response to client. size: %d", size-int(idx))
					}
					break
				}
			}
			return nil
		}

		if err != nil {
			log.Logger.Errorf("create or connect tap failed, err: (%v)", err)
		}

		if err := star.register(); err != nil {
			log.Logger.Errorf("registry failed. (%v)", err)
		}
	})
	star.starLoop()
	return nil
}
