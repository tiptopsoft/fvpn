package server

import (
	"github.com/topcloudz/fvpn/pkg/handler"
	"github.com/topcloudz/fvpn/pkg/middleware"
	"github.com/topcloudz/fvpn/pkg/middleware/infra"
	"net"
	"sync"

	"github.com/topcloudz/fvpn/pkg/cache"
	"github.com/topcloudz/fvpn/pkg/log"
	"github.com/topcloudz/fvpn/pkg/option"
	"github.com/topcloudz/fvpn/pkg/packet"
	"github.com/topcloudz/fvpn/pkg/socket"
	"golang.org/x/sys/unix"
)

var (
	once   sync.Once
	logger = log.Log()
)

// RegServer use as server
type RegServer struct {
	*option.ServerConfig
	socket   socket.Interface
	cache    *cache.Cache
	packet   packet.Interface
	ws       sync.WaitGroup
	h        handler.Handler
	Inbound  chan *packet.Frame //used from udp
	Outbound chan *packet.Frame //used for tun
}

func (r *RegServer) Start(address string) error {
	go func() {
		r.start(address)
	}()

	//启动udp处理goroutine
	go r.ReadFromUdp()
	go r.WriteToUdp()

	go func() {
		hs := New(r.cache)
		hs.Start()
	}()

	r.ws.Add(1)
	r.ws.Wait()
	return nil
}

// Peer register cache for net, and for user create client
func (r *RegServer) start(address string) error {
	r.socket = socket.NewSocket(4000)
	once.Do(func() {
		r.cache = cache.New()
		r.h = middleware.WithMiddlewares(r.serverUdpHandler(), infra.Middlewares(false, false)...)
	})

	logger.Debugf("server start at: %s", address)

	return nil
}

func (r *RegServer) initHandler() {
	r.h = middleware.WithMiddlewares(r.serverUdpHandler(), infra.Middlewares()...)
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
