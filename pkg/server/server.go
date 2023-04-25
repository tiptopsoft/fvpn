package server

import (
	"fmt"
	"github.com/topcloudz/fvpn/pkg/nativehttp"
	"net"
	"sync"

	"github.com/topcloudz/fvpn/pkg/cache"
	"github.com/topcloudz/fvpn/pkg/log"
	"github.com/topcloudz/fvpn/pkg/option"
	"github.com/topcloudz/fvpn/pkg/packet"
	"github.com/topcloudz/fvpn/pkg/packet/register"
	"github.com/topcloudz/fvpn/pkg/socket"
	"golang.org/x/sys/unix"
)

var (
	once   sync.Once
	logger = log.Log()
)

// RegStar use as server
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
		hs := nativehttp.New(r.cache)
		hs.Start()
	}()

	r.ws.Add(1)
	r.ws.Wait()
	return nil
}

// Peer register cache for net, and for user create client
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
		logger.Infof("server start at: %s", address)
		if err != nil {
			return err
		}
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
