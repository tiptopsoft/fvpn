package edge

import (
	"fmt"
	"github.com/interstellar-cloud/star/pkg/util"
	"github.com/interstellar-cloud/star/pkg/util/log"
	"github.com/interstellar-cloud/star/pkg/util/option"
	"github.com/interstellar-cloud/star/pkg/util/packet/peer"
	"github.com/interstellar-cloud/star/pkg/util/packet/register"
	"github.com/interstellar-cloud/star/pkg/util/socket"
	"github.com/interstellar-cloud/star/pkg/util/socket/executor"
	"github.com/interstellar-cloud/star/pkg/util/tuntap"
	"golang.org/x/sys/unix"
	"net"
	"os"
)

type Star struct {
	*option.StarConfig
	*tuntap.Tuntap
	currentSupernode *net.UDPConn
	Peers            map[string]Peer //获取回来的Peers  mac: Peer
	SocketExecutor   executor.Executor
	TapExecutor      executor.Executor
}

var (
	stopCh = make(chan int, 1)
)

type Peer struct {
	Conn    *net.Conn
	MacAddr net.HardwareAddr
	IP      net.IP
	Port    uint16
}

// Start logic: start to: 1. PING to registry node 2. registry to registry 3. auto ip config tuntap 4.
func (star Star) Start() error {
	//init connect to registry
	var conn net.Conn
	var err error

	conn, err = star.conn()
	if err != nil {
		return err
	}

	i := 1
loop:
	for {
		//registry to registry
		switch i {
		case 1: //registry
			err = star.register(conn)
			if err != nil {
				return err
			}
			i++
			break
		case 2: //after registry, send query
			err = star.queryPeer(conn)
			if err != nil {
				return err
			}
			i++
			break
		case 3: // start to init connect to dst
			option.AddrMap.Range(func(key, value any) bool {
				return true
			})
			i++
			break
		case 4:
			break loop
		}
	}

	//netFile, err := conn.(*net.UDPConn).File()
	tap, err := tuntap.New(tuntap.TAP)
	star.Tuntap = tap
	if err != nil {
		log.Logger.Errorf("create or connect tuntap failed. (%v)", err)
	}

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
	return conn, nil
}

func (star *Star) queryPeer(conn net.Conn) error {
	cp := peer.NewPacket()
	data, err := peer.Encode(cp)
	if err != nil {
		return err
	}

	switch star.Protocol {
	case option.UDP:
		log.Logger.Infof("Start to query star peer info, data: (%v)", data)
		if _, err := conn.(*net.UDPConn).Write(data); err != nil {
			return nil
		}
		break
	}
	return nil
}

// register register a edgestar to center.
func (star *Star) register(conn net.Conn) error {
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
		if _, err := conn.(*net.UDPConn).Write(data); err != nil {
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
	hw, _ := net.ParseMAC(star.Tuntap.MacAddr)
	rp.SrcMac = hw
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
	netFd := socket.SocketFD(star.currentSupernode)
	tapFd := int(star.Fd)
	for {
		var FdSet unix.FdSet
		var maxFd int
		if netFd > tapFd {
			maxFd = netFd
		} else {
			maxFd = tapFd
		}
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
			executor = star.TapExecutor
		}

		if FdSet.IsSet(netFd) {
			executor = star.SocketExecutor
		}

		if s.FileDescriptor != 0 {
			if err := executor.Execute(s); err != nil {
				log.Logger.Errorf("executor execute faile: (%v)", err.Error())
			}
		}

	}
}
