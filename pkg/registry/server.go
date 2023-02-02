package registry

import (
	"context"
	"fmt"
	"github.com/interstellar-cloud/star/pkg/epoller"
	"github.com/interstellar-cloud/star/pkg/handler"
	"github.com/interstellar-cloud/star/pkg/handler/auth"
	"github.com/interstellar-cloud/star/pkg/handler/encrypt"
	"github.com/interstellar-cloud/star/pkg/packet/common"
	"github.com/interstellar-cloud/star/pkg/socket"
	"github.com/interstellar-cloud/star/pkg/util/log"
	option2 "github.com/interstellar-cloud/star/pkg/util/option"
	"net"
	"sync"
)

// mac:Pub
var m sync.Map

type Node struct {
	Mac   [4]byte
	Proto option2.Protocol
	Conn  net.Conn
	Addr  *net.UDPAddr
}

//RegStar use as registry
type RegStar struct {
	*option2.RegConfig
	handler.ChainHandler
	conn net.Conn
}

func (r *RegStar) Start(address string) error {
	return r.start(address)
}

// Node register node for net, and for user create edge
func (r *RegStar) start(address string) error {
	var ctx = context.Background()
	var conn net.Conn
	r.ChainHandler = handler.NewChainHandler()
	if r.OpenAuth {
		r.AddHandler(ctx, &auth.AuthHandler{})
	}

	if r.OpenEncrypt {
		r.AddHandler(ctx, &encrypt.StarEncrypt{})
	}

	switch r.Protocol {
	case option2.UDP:
		addr, err := ResolveAddr(address)
		if err != nil {
			return err
		}

		conn, err = net.ListenUDP("udp", addr)
		r.conn = conn
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

		eventLoop, err := epoller.NewEventLoop()
		eventLoop.Protocol = r.Protocol
		if err := eventLoop.AddFd(conn); err != nil {
			log.Logger.Errorf("add fd to epoller failed. (%v)", err)
			return err
		}

		if err != nil {
			return err
		}

		eventLoop.EventLoop(r)
	default:
		log.Logger.Info("this is a tcp server")
	}

	return nil
}

func (r *RegStar) Execute(socket socket.Socket) error {
	data := make([]byte, 2048)
	_, addr, err := socket.ReadFromUdp(data)
	if err != nil {
		fmt.Println(err)
	}

	p, err := common.Decode(data)
	if err != nil {
		fmt.Println(err)
	}

	switch p.Flags {

	case option2.MsgTypeRegisterSuper:
		r.processRegister(addr, socket.UdpSocket, data, nil)
		break
	case option2.MsgTypeQueryPeer:
		r.processPeer(addr, socket.UdpSocket)
		break
	case option2.MsgTypePacket:
		r.forward(data, &p)
		break
	}

	return nil
}

func ResolveAddr(address string) (*net.UDPAddr, error) {
	return net.ResolveUDPAddr("udp", address)
}
