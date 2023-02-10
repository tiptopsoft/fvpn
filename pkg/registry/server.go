package registry

import (
	"context"
	"fmt"
	"github.com/interstellar-cloud/star/pkg/util"
	"github.com/interstellar-cloud/star/pkg/util/epoller"
	"github.com/interstellar-cloud/star/pkg/util/handler"
	"github.com/interstellar-cloud/star/pkg/util/handler/auth"
	"github.com/interstellar-cloud/star/pkg/util/handler/encrypt"
	"github.com/interstellar-cloud/star/pkg/util/log"
	option2 "github.com/interstellar-cloud/star/pkg/util/option"
	"github.com/interstellar-cloud/star/pkg/util/packet/common"
	"github.com/interstellar-cloud/star/pkg/util/socket"
	"golang.org/x/sys/unix"
	"net"
	"sync"
)

var (
	once sync.Once
)

//RegStar use as registry
type RegStar struct {
	*option2.RegConfig
	handler.ChainHandler
	socket socket.Socket
	Nodes  util.Nodes
}

func (r *RegStar) Start(address string) error {
	return r.start(address)
}

// Node register node for net, and for user create edge
func (r *RegStar) start(address string) error {
	var ctx = context.Background()
	sock := socket.NewSocket()
	r.socket = sock
	r.ChainHandler = handler.NewChainHandler()
	if r.OpenAuth {
		r.AddHandler(ctx, &auth.AuthHandler{})
	}

	if r.OpenEncrypt {
		r.AddHandler(ctx, &encrypt.StarEncrypt{})
	}

	once.Do(func() {
		r.Nodes = make(util.Nodes)
	})

	switch r.Protocol {
	case option2.UDP:
		addr, err := ResolveAddr(address)

		if err != nil {
			return err
		}

		err = r.socket.Listen(addr)
		if err != nil {
			return err
		}
		log.Logger.Infof("registry start at: %s", address)

		//start http
		rs := RegistryServer{
			RegStar: r,
		}
		go func() {
			if err := rs.Start(); err != nil {
				log.Logger.Errorf("this is udp server, listen http failed.")
			}
		}()

		if err != nil {
			return err
		}
		//defer conn.Close()

		eventLoop, err := epoller.NewEventLoop()
		eventLoop.Protocol = r.Protocol
		if err := eventLoop.AddFd(r.socket); err != nil {
			log.Logger.Errorf("add fd to epoller failed. err: (%v)", err)
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
	size, addr, err := socket.ReadFromUdp(data)
	if err != nil {
		fmt.Println(err)
	}

	p, err := common.Decode(data)
	if err != nil {
		fmt.Println(err)
	}

	switch p.Flags {

	case option2.MsgTypeRegisterSuper:
		r.processRegister(addr, data[:size], nil)
		break
	case option2.MsgTypeQueryPeer:
		r.processFindPeer(addr)
		break
	case option2.MsgTypePacket:
		r.forward(data[:size], &p)
		break
	}

	return nil
}

func ResolveAddr(address string) (unix.Sockaddr, error) {
	addr, err := net.ResolveUDPAddr("udp", address)
	if err != nil {
		return nil, err
	}

	ip := [4]byte{}
	copy(ip[:], addr.IP.To4())

	result := &unix.SockaddrInet4{
		Port: addr.Port,
		Addr: ip,
	}

	return result, nil
}
