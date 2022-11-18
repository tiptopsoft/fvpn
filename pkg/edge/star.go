package edge

import (
	"github.com/interstellar-cloud/star/pkg/device"
	"github.com/interstellar-cloud/star/pkg/log"
	"github.com/interstellar-cloud/star/pkg/option"
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
1. PING to register node 2. register to register 3. auto ip config tuntap 4.
*/
func (edge EdgeStar) Start() error {
	//init connect to registry
	var conn net.Conn
	var err error

	conn, err = edge.conn(edge.Registry)
	if err != nil {
		return err
	}

	ch <- 1
	// registry to registry
	switch <-ch {
	case 1: //register
		err = edge.register(conn)
		if err != nil {
			return err
		}
		break
	case 2: //after register, send query
		edge.queryPeer(conn)
	}

	err = edge.process(conn)
	if err != nil {
		log.Logger.Errorf("process failed, err:%v", err)
		// re start a goroutine.
	}

	//run loop process udp
	//time.Sleep(1000 * 3)
	//go func() {
	//	err = edge.process(conn)
	//	if err != nil {
	//		log.Logger.Errorf("process failed, err:%v", err)
	//		// re start a goroutine.
	//	}
	//}()

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

	log.Logger.Info("star connected to server:", es.Registry)
	return conn, nil
}
