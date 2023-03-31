package fvpns

import (
	"fmt"
	"github.com/interstellar-cloud/star/pkg/http"
	"net"
	"sync"

	"github.com/interstellar-cloud/star/pkg/cache"
	"github.com/interstellar-cloud/star/pkg/epoller"
	"github.com/interstellar-cloud/star/pkg/log"
	"github.com/interstellar-cloud/star/pkg/option"
	"github.com/interstellar-cloud/star/pkg/packet"
	"github.com/interstellar-cloud/star/pkg/packet/register"
	"github.com/interstellar-cloud/star/pkg/socket"
	"golang.org/x/sys/unix"
)

var (
	once   sync.Once
	logger = log.Log()
)

// RegStar use as fvpns
type RegStar struct {
	*option.ServerConfig
	socket socket.Interface
	cache  cache.PeersCache
	packet packet.Interface
	ws     sync.WaitGroup
}

func (r *RegStar) Cache() cache.PeersCache {
	return r.cache
}

func (r *RegStar) Start(address string) error {
	go func() {
		r.start(address)
	}()

	go func() {
		hs := http.New(r.cache)
		hs.Start()
	}()

	r.ws.Add(1)
	r.ws.Wait()
	return nil
}

// Peer register cache for net, and for user create fvpnc
func (r *RegStar) start(address string) error {
	r.socket = socket.NewSocket()
	once.Do(func() {
		r.cache = cache.New()
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
		logger.Infof("fvpns start at: %s", address)
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

func (r *RegStar) Execute(socket socket.Interface) error {
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
