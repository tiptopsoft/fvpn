package client

import (
	"github.com/topcloudz/fvpn/pkg/middleware/encrypt"
	"sync"

	"github.com/topcloudz/fvpn/pkg/cache"
	udphandler "github.com/topcloudz/fvpn/pkg/handler/udp"
	"github.com/topcloudz/fvpn/pkg/middleware"
	"github.com/topcloudz/fvpn/pkg/middleware/auth"
	"github.com/topcloudz/fvpn/pkg/option"
	processorudp "github.com/topcloudz/fvpn/pkg/processor/udp"
	"github.com/topcloudz/fvpn/pkg/socket"
)

var (
	once        sync.Once
	DefaultPort = 6663
)

type Node struct {
	*option.Config
	Protocol  option.Protocol
	tuns      sync.Map //key: netId, value: Tuntap
	socket    socket.Interface
	cache     cache.PeersCache //获取回来的Peers  mac: Peer
	processor sync.Map         //核心处理逻辑
}

func (n *Node) Start() error {
	once.Do(func() {
		n.socket = socket.NewSocket()
		if err := n.conn(); err != nil {
			logger.Errorf("failed to connect to server: %v", err)
		}
		n.cache = cache.New()
		n.Protocol = option.UDP
		//n.initExecutor()
		n.initUdpHandler()
	})
	go n.starLoop()
	return n.runHttpServer()
}

func (n *Node) initUdpHandler() {
	udpHandler := middleware.WithMiddlewares(udphandler.New(&n.tuns, n.cache), n.initMiddleware()...)
	udpProcessor := processorudp.New(udpHandler, n.socket)
	n.processor.Store(n.socket.(socket.Socket).Fd, udpProcessor)
}

func (n *Node) initMiddleware() []middleware.Middleware {
	var result []middleware.Middleware
	if n.OpenAuth {
		result = append(result, auth.Middleware())
	}

	if n.OpenEncrypt {
		result = append(result, encrypt.Middleeare())
	}

	if n.OpenCompress {
		//TODO
	}

	return result
}
