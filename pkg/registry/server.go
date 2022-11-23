package registry

import (
	"context"
	"fmt"
	"github.com/interstellar-cloud/star/pkg/handler"
	"github.com/interstellar-cloud/star/pkg/handler/auth"
	"github.com/interstellar-cloud/star/pkg/handler/encrypt"
	"github.com/interstellar-cloud/star/pkg/log"
	"github.com/interstellar-cloud/star/pkg/option"
	"github.com/interstellar-cloud/star/pkg/packet/common"
	"net"
	"sync"
)

var limitChan = make(chan int, 1)

// mac:Pub
var m sync.Map

type Node struct {
	Mac   [4]byte
	Proto option.Protocol
	Conn  net.Conn
	Addr  *net.UDPAddr
}

//RegStar use as registry
type RegStar struct {
	*option.RegConfig
	handler.Executor
	conn net.Conn
}

func (r *RegStar) Start(address string) error {
	return r.start(address)
}

// Node register node for net, and for user create edge
func (r *RegStar) start(address string) error {
	var ctx = context.Background()
	var conn net.Conn
	r.Executor = handler.NewExecutor()
	if r.OpenAuth {
		r.AddHandler(ctx, &auth.AuthHandler{})
	}

	if r.OpenEncrypt {
		r.AddHandler(ctx, &encrypt.StarEncrypt{})
	}

	switch r.Protocol {
	case option.UDP:
		addr, err := ResolveAddr(address)
		if err != nil {
			return err
		}

		conn, err = net.ListenUDP("udp", addr)

		log.Logger.Infof("registry start at: %s", address)

		//start http
		rs := RegistryServer{
			r.RegConfig,
		}
		go func() {
			if err := rs.Start(rs.HttpListen); err != nil {
				log.Logger.Errorf("this is udp server, listen http failed.")
			}
		}()

		if err != nil {
			return err
		}
		defer conn.Close()
		for {
			limitChan <- 1
			go r.handleUdp(ctx, conn.(*net.UDPConn))
		}
	default:
		log.Logger.Info("this is a tcp server")
	}

	return nil
}

func (r *RegStar) handleUdp(ctx context.Context, conn *net.UDPConn) {
	for {
		data := make([]byte, 2048)
		_, addr, err := conn.ReadFromUDP(data)
		if err != nil {
			fmt.Println(err)
		}

		p, err := common.Decode(data)
		if err != nil {
			fmt.Println(err)
		}

		//exec executor
		if err := r.Execute(ctx, data); err != nil {
			fmt.Println(err)
		}

		switch p.Flags {

		case option.MSG_TYPE_REGISTER_SUPER:
			r.processRegister(addr, conn, data, nil)
			break
		case option.MSG_TYPE_QUERY_PEER:
			r.processPeer(addr, conn, data, &p)
			break
		case option.MSG_TYPE_PACKET:
			r.forward(data, &p)
			break
		}
	}

}

func ResolveAddr(address string) (*net.UDPAddr, error) {
	return net.ResolveUDPAddr("udp", address)
}
