package relay

import (
	"github.com/topcloudz/fvpn/pkg/util"
	"net"
	"sync"

	"github.com/topcloudz/fvpn/pkg/log"
	"github.com/topcloudz/fvpn/pkg/packet"
)

var (
	once   sync.Once
	logger = log.Log()
)

// RegServer use as server
type RegServer struct {
	*util.ServerConfig
	//socket   socket.Interface
	socket *net.UDPConn
	ws     sync.WaitGroup
	//readHandler  handler.Handler
	//writeHandler handler.Handler
	Inbound  chan *packet.Frame //used from udp
	Outbound chan *packet.Frame //used for tun

	//every node has it's own key
	manager *util.KeyManager
	appIds  map[string]string
}

func (r *RegServer) Start(address string) error {
	r.manager = &util.KeyManager{NodeKeys: make(map[string]*util.NodeKey, 1)}
	err := r.start(address)
	if err != nil {
		return err
	}

	r.ws.Add(1)
	r.ws.Wait()
	return nil
}

// Peer register cache for net, and for user create client
func (r *RegServer) start(address string) error {
	//r.socket = socket.NewSocket(4000)
	socket, _ := net.ListenUDP("udp", &net.UDPAddr{
		IP: net.IPv4zero, Port: 4000})
	//once.Do(func() {
	//	r.cache = cache.New()
	//	r.readHandler = middleware.WithMiddlewares(r.serverUdpHandler(), device.Decode(r.manager))
	//	r.writeHandler = middleware.WithMiddlewares(r.writeUdpHandler(), device.Encode(r.manager))
	//})
	r.socket = socket

	logger.Debugf("server start at: %s", address)

	return nil
}
