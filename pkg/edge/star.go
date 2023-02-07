package edge

import (
	"fmt"
	"github.com/interstellar-cloud/star/pkg/util"
	"github.com/interstellar-cloud/star/pkg/util/log"
	"github.com/interstellar-cloud/star/pkg/util/option"
	"github.com/interstellar-cloud/star/pkg/util/packet"
	"github.com/interstellar-cloud/star/pkg/util/packet/common"
	"github.com/interstellar-cloud/star/pkg/util/packet/forward"
	"github.com/interstellar-cloud/star/pkg/util/packet/peer"
	peerack "github.com/interstellar-cloud/star/pkg/util/packet/peer/ack"
	"github.com/interstellar-cloud/star/pkg/util/packet/register"
	"github.com/interstellar-cloud/star/pkg/util/packet/register/ack"
	"github.com/interstellar-cloud/star/pkg/util/socket"
	"github.com/interstellar-cloud/star/pkg/util/socket/executor"
	"github.com/interstellar-cloud/star/pkg/util/tuntap"
	"golang.org/x/sys/unix"
	"io"
	"net"
	"os"
	"sync"
	"time"
	"unsafe"
)

type Star struct {
	*option.StarConfig
	*tuntap.Tuntap
	socket.Socket
	Peers          util.Peers //获取回来的Peers  mac: Peer
	Executor       executor.Executor
	SocketExecutor executor.Executor
	TapExecutor    executor.Executor
}

type TapExecutor struct {
	Tap       *tuntap.Tuntap
	NetSocket socket.Socket
	TapSocket socket.Socket
	Protocol  option.Protocol
	Peers     util.Peers
}

type NetExecutor struct {
	Tap       *tuntap.Tuntap
	NetSocket socket.Socket
	TapSocket socket.Socket
	Protocol  option.Protocol
	Peers     util.Peers //获取回来的Peers  mac: Peer
}

var (
	stopCh = make(chan int, 1)
	once   sync.Once
)

// Start logic: start to: 1. PING to registry node 2. registry to registry 3. auto ip config tuntap 4.
func (star Star) Start() error {
	//init connect to registry
	var err error
	_, err = star.conn()
	if err != nil {
		return err
	}
	err = star.register()

	go func() {
		for {
			err = star.queryPeer()
			time.Sleep(time.Second * 60)
		}
	}()

	once.Do(func() {
		star.Peers = make(util.Peers, 1)
		tap, err := tuntap.New(tuntap.TAP)
		star.Tuntap = tap
		tapSocket := socket.Socket{
			AppType:        option.UDP,
			FileDescriptor: int(star.Tuntap.Fd),
			UdpSocket:      star.UdpSocket,
		}
		star.SocketExecutor = NetExecutor{
			Tap:       tap,
			NetSocket: star.Socket,
			TapSocket: tapSocket,
			Protocol:  star.Protocol,
			Peers:     star.Peers,
		}

		star.TapExecutor = TapExecutor{
			Tap:       tap,
			NetSocket: star.Socket,
			TapSocket: tapSocket,
			Protocol:  star.Protocol,
			Peers:     star.Peers,
		}
		if err != nil {
			log.Logger.Errorf("create or connect tuntap failed. (%v)", err)
		}
	})

	star.loop()
	if <-stopCh > 0 {
		log.Logger.Infof("star stop success")
		os.Exit(-1)
	}
	return nil
}

func (star *Star) conn() (net.Conn, error) {
	var conn net.Conn
	var err error

	switch star.Protocol {
	case option.UDP:
		conn, err = net.Dial("udp", star.Registry)
	}

	//defer conn.Close()
	if err != nil {
		return nil, err
	}

	log.Logger.Infof("star connected to registry: (%v)", star.Registry)
	star.Socket = socket.Socket{
		AppType:   option.UDP,
		UdpSocket: conn.(*net.UDPConn),
	}
	return conn, nil
}

func (star *Star) queryPeer() error {
	cp := peer.NewPacket()
	data, err := peer.Encode(cp)
	if err != nil {
		return err
	}

	switch star.Protocol {
	case option.UDP:
		log.Logger.Infof("start to query star peer info, data: (%v)", data)
		if _, err := star.Write(data); err != nil {
			return nil
		}
		break
	}
	return nil
}

// register register a edgestar to center.
func (star *Star) register() error {
	var err error
	rp := register.NewPacket()
	hw, _ := net.ParseMAC(util.GetLocalMacAddr())
	rp.SrcMac = hw
	data, err := register.Encode(rp)
	log.Logger.Infof("sending registry data: %v", data)
	if err != nil {
		return err
	}

	switch star.Protocol {
	case option.UDP:
		log.Logger.Infof("star start to registry self to registry: %v", rp)
		if _, err := star.Write(data); err != nil {
			return err
		}
		break
	}
	return nil
}

// register register a edgestar to center.
func (star *Star) unregister(conn net.Conn) error {
	var err error

	rp := register.NewUnregisterPacket()
	rp.SrcMac = star.Tuntap.MacAddr
	data, err := register.Encode(rp)
	fmt.Println("sending unregister data: ", data)
	if err != nil {
		return err
	}

	switch star.Protocol {
	case option.UDP:
		log.Logger.Infof("star start to registry self to registry: %v", rp)
		if _, err := conn.(*net.UDPConn).Write(data); err != nil {
			return err
		}
		break
	}
	return nil
}

func (star *Star) loop() {
	netFd := socket.SocketFD(star.UdpSocket)
	tapFd := int(star.Fd)
	var FdSet unix.FdSet
	var maxFd int
	if netFd > tapFd {
		maxFd = netFd
	} else {
		maxFd = tapFd
	}
	for {
		FdSet.Zero()
		FdSet.Set(netFd)
		FdSet.Set(tapFd)

		ret, err := unix.Select(maxFd+1, &FdSet, nil, nil, nil)
		if ret < 0 && err == unix.EINTR {
			continue
		}
		var s socket.Socket
		var executor executor.Executor
		if err != nil {
			panic(err)
		}

		if FdSet.IsSet(tapFd) {
			s = socket.Socket{
				FileDescriptor: tapFd,
				UdpSocket:      nil,
			}
			executor = star.TapExecutor
		}

		if FdSet.IsSet(netFd) {
			s = socket.Socket{
				FileDescriptor: netFd,
				UdpSocket:      star.UdpSocket,
			}
			executor = star.SocketExecutor
		}

		if err := executor.Execute(s); err != nil {
			log.Logger.Errorf("executor execute faile: (%v)", err.Error())
		}

	}
}

func (ee NetExecutor) Execute(socket socket.Socket) error {
	log.Logger.Infof("start execute net...")
	if ee.Protocol == option.UDP {
		//for {
		udpBytes := make([]byte, 2048)
		n, err := socket.Read(udpBytes)
		log.Logger.Infof("star net socket receive size: %d, data: (%v)", n, udpBytes)
		if err != nil {
			if err == io.EOF {
				//no data exists, continue read next frame continue
				log.Logger.Errorf("not data exists")
			} else {
				log.Logger.Errorf("read from remote error: %v", err)
			}
		}

		cp, err := common.Decode(udpBytes)

		log.Logger.Infof("edge net executor working...., data: %v", cp)
		if err != nil {
			log.Logger.Errorf("decode err: %v", err)
		}

		switch cp.Flags {
		case option.MsgTypeRegisterAck:
			regAck, err := ack.Decode(udpBytes)
			if err != nil {
				return err
			}
			log.Logger.Infof("got registry registry ack: %v", regAck)
			//设置IP
			if err = option.ExecCommand("/bin/sh", "-c", fmt.Sprintf("ifconfig %s %s netmask %s mtu %d up", ee.Tap.Name, regAck.AutoIP.String(), regAck.Mask.String(), 1420)); err != nil {
				return err
			}
			break
		case option.MsgTypeQueryPeer:
			fmt.Println("got registry peers start...")
			//get peerInfo
			peerPacketAck, err := peerack.Decode(udpBytes)
			if err != nil {
				return err
			}
			infos := peerPacketAck.PeerInfos
			log.Logger.Infof("got registry peers: (%v)", infos)
			for _, v := range infos {
				addr, err := net.ResolveUDPAddr("udp", fmt.Sprintf("%s:%d", v.Host.String(), v.Port))
				if err != nil {
					log.Logger.Errorf("resolve addr failed. err: %v", err)
				}
				//option.AddrMap.Store(v.Mac.String(), addr)
				conn, err := net.Dial("udp", addr.String())
				if err != nil {
					return err
				}
				peer := &util.Peer{
					Conn:    conn,
					MacAddr: v.Mac,
					IP:      v.Host,
					Port:    v.Port,
				}
				ee.Peers[v.Mac.String()] = peer
			}
			break
		case option.MsgTypePacket:
			forwardPacket, err := forward.Decode(udpBytes)
			if err != nil {
				return err
			}
			log.Logger.Infof("got through packet: %v, srcMac: %v, current tap macAddr: %v", forwardPacket, forwardPacket.SrcMac, ee.Tap.MacAddr)

			if forwardPacket.SrcMac.String() == ee.Tap.MacAddr.String() {
				//self, drop packet
				log.Logger.Warnf("self packet droped: %v, srcMac: %v, current tap macAddr: %v", forwardPacket, forwardPacket.SrcMac, ee.Tap.MacAddr)
			} else {
				//写入到tap
				idx := unsafe.Sizeof(forwardPacket)
				if _, err := ee.Tap.Socket.Write(udpBytes[idx:n]); err != nil {
					log.Logger.Errorf("write to tap failed. (%v)", err.Error())
				}
				log.Logger.Infof("net write to tap as tap response to client. size: %d", n-int(idx))
			}
			break
		}
	}
	return nil
}

// Execute TapExecutor use to handle tap frame, write to udp sock.
// Read a single packet from the TAP interface, process it and write out the corresponding packet to the cooked socket.
func (te TapExecutor) Execute(socket socket.Socket) error {
	b := make([]byte, option.STAR_PKT_BUFF_SIZE)
	n, err := socket.Read(b)
	log.Logger.Info(fmt.Sprintf("Read from tap %s: length: %d，data: %v", te.Tap.Name, n, b))
	if err != nil {
		log.Logger.Errorf("tap read failed. (%v)", err)
		return err
	}

	destMac := getMacAddr(b)
	log.Logger.Infof("Tap dev: %s receive: %d byte, mac: %v", te.Tap.Name, n, destMac)
	//broad := util.IsBroadCast(destMac)
	//if broad {
	// broad frame, go through supernode
	fp := forward.NewPacket()
	fp.SrcMac, err = util.GetMacAddrByDev(te.Tap.Name)
	if err != nil {
		log.Logger.Errorf("get src mac failed. %v", err)
	}
	fp.DstMac, err = net.ParseMAC(destMac)
	if err != nil {
		log.Logger.Errorf("get src mac failed. %v", err)
	}

	bs, err := forward.Encode(fp)
	if err != nil {
		log.Logger.Errorf("encode forward failed. err: %v", err)
	}

	idx := 0
	newPacket := make([]byte, 2048)
	idx = packet.EncodeBytes(newPacket, bs, idx)
	packet.EncodeBytes(newPacket, b[:n], idx)
	write2Net(te.NetSocket, newPacket)
	//} else {
	//	// go p2p
	//}
	return nil
}

//use host socket write so destination
func write2Net(socket socket.Socket, b []byte) {
	log.Logger.Infof("tap write to net packet: (%v)", b)
	if _, err := socket.Write(b); err != nil {
		log.Logger.Errorf("write to remote failed. (%v)", err)
	}
}

func getMacAddr(buf []byte) string {
	return fmt.Sprintf("%02x:%02x:%02x:%02x:%02x:%02x", buf[0], buf[1], buf[2], buf[3], buf[4], buf[5])
}
