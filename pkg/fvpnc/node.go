package fvpnc

import (
	"github.com/interstellar-cloud/star/pkg/middleware/encrypt"
	"sync"
	"time"

	"github.com/interstellar-cloud/star/pkg/handler/device"
	"github.com/interstellar-cloud/star/pkg/handler/udp"

	"github.com/interstellar-cloud/star/pkg/cache"
	"github.com/interstellar-cloud/star/pkg/middleware"
	"github.com/interstellar-cloud/star/pkg/middleware/auth"
	"github.com/interstellar-cloud/star/pkg/option"
	processordevice "github.com/interstellar-cloud/star/pkg/processor/device"
	processorudp "github.com/interstellar-cloud/star/pkg/processor/udp"
	"github.com/interstellar-cloud/star/pkg/socket"
	"github.com/interstellar-cloud/star/pkg/tuntap"
)

var (
	once            sync.Once
	DefaultEdgePort = 6061
)

type Node struct {
	*option.Config
	Protocol  option.Protocol
	tap       *tuntap.Tuntap
	socket    socket.Interface
	cache     cache.PeersCache //获取回来的Peers  mac: Peer
	processor sync.Map         //核心处理逻辑
}

func (node *Node) Start() error {
	once.Do(func() {
		node.socket = socket.NewSocket()
		if err := node.conn(); err != nil {
			logger.Errorf("failed to connect to fvpns: %v", err)
		}
		node.cache = cache.New()
		node.Protocol = option.UDP
		tap, err := tuntap.New(tuntap.TAP)
		node.tap = tap

		if err != nil {
			logger.Errorf("create or connect tap failed, err: (%v)", err)
		}

		if err := node.register(); err != nil {
			logger.Errorf("fvpns failed. (%v)", err)
		}

		//node.initExecutor()
		node.initProcessor()
		go func() {
			for {
				node.queryPeer()
				//连通
				node.dialNode()
				time.Sleep(30 * time.Second)
			}
		}()
	})
	node.starLoop()
	return nil
}

func (node *Node) initProcessor() {
	deviceHandler := middleware.WithMiddlewares(device.New(node.tap, node.socket, node.cache), node.initMiddleware()...)
	deviceProcessor := processordevice.New(node.tap, deviceHandler)
	udpHandler := middleware.WithMiddlewares(udp.New(node.tap, node.cache), node.initMiddleware()...)
	udpProcessor := processorudp.New(node.tap, udpHandler, node.socket)

	node.processor.Store(node.tap.Fd, deviceProcessor)
	node.processor.Store(node.socket.(socket.Socket).Fd, udpProcessor)
}

func (node *Node) initMiddleware() []middleware.Middleware {
	var result []middleware.Middleware
	if node.OpenAuth {
		result = append(result, auth.Middleware())
	}

	if node.OpenEncrypt {
		result = append(result, encrypt.Middleeare())
	}

	if node.OpenCompress {
		//TODO
	}

	return node.initMiddleware()
}
