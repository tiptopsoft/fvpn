package edge

import (
	"github.com/interstellar-cloud/star/pkg/log"
	"github.com/interstellar-cloud/star/pkg/option"
	"github.com/interstellar-cloud/star/pkg/socket"
	"github.com/interstellar-cloud/star/pkg/tuntap"
	"net"
	"os"
)

type StarEdge struct {
	*option.EdgeConfig
	tap *tuntap.Tuntap
}

var (
	stopCh = make(chan int, 1)
	ch     = make(chan int, 1)
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

	ch <- 1
	// registry to registry
	switch <-ch {
	case 1: //registry
		err = edge.register(conn)
		if err != nil {
			return err
		}
		break
	case 2: //after registry, send query
		err = edge.queryPeer(conn)
		if err != nil {
			return err
		}
		break
	case 3: // start to init connect to dst
		option.AddrMap.Range(func(key, value any) bool {
			return true
		})
		break
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
