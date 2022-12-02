package edge

import (
	"github.com/interstellar-cloud/star/pkg/device"
	"github.com/interstellar-cloud/star/pkg/log"
	"github.com/interstellar-cloud/star/pkg/option"
	"github.com/interstellar-cloud/star/pkg/socket"
	"net"
	"os"
)

type EdgeStar struct {
	*option.EdgeConfig
	tap *device.Tuntap
}

var (
	stopCh = make(chan int, 1)
	ch     = make(chan int, 1)
)

/**
 * Start logic: start to:
1. PING to registry node 2. registry to registry 3. auto ip config tuntap 4.
*/
func (edge EdgeStar) Start() error {
	//init connect to registry
	var conn net.Conn
	var err error

	conn, err = edge.conn(edge.Registry)
	if err != nil {
		return err
	}

	s, err := socket.NewSocket(conn.(*net.UDPConn))
	if err != nil {
		return err
	}

	eventLoop, err := socket.NewEventLoop(s)
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

	eventLoop.EventLoop()

	if <-stopCh > 0 {
		log.Logger.Infof("edge stop success")
		os.Exit(-1)
	}
	return nil
}

func (es *EdgeStar) conn(address string) (net.Conn, error) {
	var conn net.Conn
	var err error

	switch es.Protocol {
	case option.UDP:
		conn, err = net.Dial("udp", es.Registry)
	}

	//defer conn.Close()
	if err != nil {
		return nil, err
	}

	log.Logger.Info("star connected to registry:", es.Registry)
	return conn, nil
}
