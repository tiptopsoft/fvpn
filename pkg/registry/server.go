package registry

import (
	"fmt"
	"github.com/interstellar-cloud/star/pkg/epoller"
	"github.com/interstellar-cloud/star/pkg/handler"
	"github.com/interstellar-cloud/star/pkg/log"
	"github.com/interstellar-cloud/star/pkg/node"
	"github.com/interstellar-cloud/star/pkg/option"
	"github.com/interstellar-cloud/star/pkg/packet"
	"github.com/interstellar-cloud/star/pkg/packet/register"
	socket2 "github.com/interstellar-cloud/star/pkg/socket"
	"golang.org/x/sys/unix"
	"net"
	"sync"
)

var (
	once   sync.Once
	logger = log.Log()
)

//RegStar use as registry
type RegStar struct {
	*option.RegConfig
	socket      socket2.Interface
	cache       node.NodesCache
	AuthHandler handler.Interface
	packet      packet.Interface
}

func (r *RegStar) Start(address string) error {
	return r.start(address)
}

// Node register node for net, and for user create edge
func (r *RegStar) start(address string) error {
	r.socket = socket2.NewSocket()
	once.Do(func() {
		r.cache = node.New()
	})

	switch r.Protocol {
	case option.UDP:
		addr, err := ResolveAddr(address)
		if err != nil {
			return err
		}

		err = r.socket.Listen(addr)
		if err != nil {
			return err
		}
		logger.Infof("registry start at: %s", address)
		if err != nil {
			return err
		}

		eventLoop, err := epoller.NewEventLoop()
		eventLoop.Protocol = r.Protocol
		if err := eventLoop.AddFd(r.socket); err != nil {
			logger.Errorf("add fd to epoller failed. err: (%v)", err)
			return err
		}

		if err != nil {
			return err
		}

		eventLoop.EventLoop(r)
	default:
		logger.Info("this is a tcp server")
	}

	return nil
}

func (r *RegStar) Execute(socket socket2.Interface) error {
	data := make([]byte, 2048)
	size, addr, err := socket.ReadFromUdp(data)
	if err != nil {
		fmt.Println(err)
	}

	pInterface, err := packet.NewPacketWithoutType().Decode(data)
	p := pInterface.(packet.Header)

	if err != nil {
		fmt.Println(err)
	}

	switch p.Flags {

	case option.MsgTypeRegisterSuper:
		r.packet = register.NewPacket()
		r.processRegister(addr, data[:size], nil)
		break
	case option.MsgTypeQueryPeer:
		r.processFindPeer(addr)
		break
	case option.MsgTypePacket:
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
