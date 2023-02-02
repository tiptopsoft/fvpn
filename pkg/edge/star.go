package edge

import (
	"fmt"
	"github.com/interstellar-cloud/star/pkg/packet/peer"
	"github.com/interstellar-cloud/star/pkg/packet/register"
	"github.com/interstellar-cloud/star/pkg/socket"
	"github.com/interstellar-cloud/star/pkg/tuntap"
	"github.com/interstellar-cloud/star/pkg/util/log"
	"github.com/interstellar-cloud/star/pkg/util/option"
	"net"
	"os"
)

type StarEdge struct {
	*option.EdgeConfig
	tap *tuntap.Tuntap
}

var (
	stopCh = make(chan int, 1)
	//ch     = make(chan int, 1)
)

// Start logic: start to: 1. PING to registry node 2. registry to registry 3. auto ip config tuntap 4.
func (edge StarEdge) Start() error {
	//init connect to registry
	var conn net.Conn
	var err error

	conn, err = edge.conn()
	if err != nil {
		return err
	}

	i := 1
outloop:
	for {
		//registry to registry
		switch i {
		case 1: //registry
			err = edge.register(conn)
			if err != nil {
				return err
			}
			i++
			break
		case 2: //after registry, send query
			err = edge.queryPeer(conn)
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
			break outloop
		}
	}

	//netFile, err := conn.(*net.UDPConn).File()
	netFd := socket.SocketFD(conn)
	tap, err := tuntap.New(tuntap.TAP)
	if err != nil {
		log.Logger.Errorf("create or connect tuntap failed. (%v)", err)
	}

	eventLoop := EventLoop{Tap: tap}
	eventLoop.eventLoop(netFd, int(tap.Fd))

	if <-stopCh > 0 {
		log.Logger.Infof("edge stop success")
		os.Exit(-1)
	}
	return nil
}

func (edge *StarEdge) conn() (net.Conn, error) {
	var conn net.Conn
	var err error

	switch edge.Protocol {
	case option.UDP:
		conn, err = net.Dial("udp", edge.Registry)
	}

	//defer conn.Close()
	if err != nil {
		return nil, err
	}

	log.Logger.Infof("star connected to registry: (%v)", edge.Registry)
	return conn, nil
}

func (edge *StarEdge) queryPeer(conn net.Conn) error {
	cp := peer.NewPacket()
	data, err := peer.Encode(cp)
	if err != nil {
		return err
	}

	switch edge.Protocol {
	case option.UDP:
		log.Logger.Infof("Start to query edge peer info, data: (%v)", data)
		if _, err := conn.(*net.UDPConn).Write(data); err != nil {
			return nil
		}
		break
	}
	return nil
}

// register register a edgestar to center.
func (edge *StarEdge) register(conn net.Conn) error {
	var err error
	rp := register.NewPacket()
	hw, _ := net.ParseMAC(GetLocalMacAddr())
	rp.SrcMac = hw
	data, err := register.Encode(rp)
	log.Logger.Infof("sending registry data: %v", data)
	if err != nil {
		return err
	}

	switch edge.Protocol {
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
func (edge *StarEdge) unregister(conn net.Conn) error {
	var err error

	rp := register.NewUnregisterPacket()
	hw, _ := net.ParseMAC(edge.MacAddr)
	rp.SrcMac = hw
	data, err := register.Encode(rp)
	fmt.Println("sending unregister data: ", data)
	if err != nil {
		return err
	}

	switch edge.Protocol {
	case option.UDP:
		log.Logger.Infof("star start to registry self to registry: %v", rp)
		if _, err := conn.(*net.UDPConn).Write(data); err != nil {
			return err
		}
		break
	}
	return nil
}
