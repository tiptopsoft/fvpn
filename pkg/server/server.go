package server

import (
	"github.com/topcloudz/fvpn/pkg/handler"
	"github.com/topcloudz/fvpn/pkg/middleware"
	"github.com/topcloudz/fvpn/pkg/middleware/codec"
	"github.com/topcloudz/fvpn/pkg/security"
	"net"
	"sync"

	"github.com/topcloudz/fvpn/pkg/cache"
	"github.com/topcloudz/fvpn/pkg/log"
	"github.com/topcloudz/fvpn/pkg/option"
	"github.com/topcloudz/fvpn/pkg/packet"
)

var (
	once   sync.Once
	logger = log.Log()
)

// RegServer use as server
type RegServer struct {
	*option.ServerConfig
	//socket   socket.Interface
	socket       *net.UDPConn
	cache        *cache.Cache
	packet       packet.Interface
	ws           sync.WaitGroup
	readHandler  handler.Handler
	writeHandler handler.Handler
	Inbound      chan *packet.Frame //used from udp
	Outbound     chan *packet.Frame //used for tun
	PrivateKey   security.NoisePrivateKey
	PubKey       security.NoisePublicKey
	cipher       security.CipherFunc
}

func (r *RegServer) Start(address string) error {
	err := r.start(address)
	if err != nil {
		return err
	}

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
	//r.socket = socket.NewSocket(4000)
	socket, _ := net.ListenUDP("udp", &net.UDPAddr{
		IP: net.IPv4zero, Port: 4000})
	once.Do(func() {
		r.cache = cache.New()
		r.readHandler = middleware.WithMiddlewares(r.serverUdpHandler(), codec.Decode(r.cipher))
		r.writeHandler = middleware.WithMiddlewares(r.writeUdpHandler(), codec.Encode(r.cipher))
	})
	r.socket = socket

	logger.Debugf("server start at: %s", address)

	return nil
}
