package client

import (
	encrypt "github.com/interstellar-cloud/star/pkg/middleware/decrypt"
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
	*option.StarConfig
	tap       *tuntap.Tuntap
	socket    socket.Interface
	cache     cache.PeersCache //获取回来的Peers  mac: Peer
	processor sync.Map         //核心处理逻辑
}

func (node *Node) Start() error {
	once.Do(func() {
		node.socket = socket.NewSocket()
		if err := node.conn(); err != nil {
			logger.Errorf("failed to connect to registry: %v", err)
		}
		node.cache = cache.New()
		node.Protocol = option.UDP
		tap, err := tuntap.New(tuntap.TAP)
		node.tap = tap

		if err != nil {
			logger.Errorf("create or connect tap failed, err: (%v)", err)
		}

		if err := node.register(); err != nil {
			logger.Errorf("registry failed. (%v)", err)
		}

		//fvpn.initExecutor()
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
	deviceHandler := middleware.WithMiddlewares(device.New(node.tap, node.socket, node.cache), node.initMiddlewares()...)
	deviceProcessor := processordevice.New(node.tap, deviceHandler)
	udpHandler := middleware.WithMiddlewares(udp.New(node.tap, node.cache), node.initMiddlewares()...)
	udpProcessor := processorudp.New(node.tap, udpHandler, node.socket)

	node.processor.Store(node.tap.Fd, deviceProcessor)
	node.processor.Store(node.socket.(socket.Socket).Fd, udpProcessor)
}

func (node *Node) initMiddlewares() []middleware.Middleware {
	cfg := node.StarConfig
	var res []middleware.Middleware
	if cfg.OpenAuth {
		res = append(res, auth.Middleware())
	}

	if cfg.OpenEncrypt {
		res = append(res, encrypt.Middleeare())
	}

	return res
}
