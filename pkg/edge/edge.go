package edge

import (
	"github.com/interstellar-cloud/star/pkg/handler/device"
	"github.com/interstellar-cloud/star/pkg/handler/udp"
	"sync"
	"time"

	"github.com/interstellar-cloud/star/pkg/executor"
	"github.com/interstellar-cloud/star/pkg/middleware"
	"github.com/interstellar-cloud/star/pkg/middleware/auth"
	"github.com/interstellar-cloud/star/pkg/node"
	"github.com/interstellar-cloud/star/pkg/option"
	"github.com/interstellar-cloud/star/pkg/packet"
	processordevice "github.com/interstellar-cloud/star/pkg/processor/device"
	processorudp "github.com/interstellar-cloud/star/pkg/processor/udp"
	"github.com/interstellar-cloud/star/pkg/socket"
	"github.com/interstellar-cloud/star/pkg/tuntap"
)

var (
	once            sync.Once
	DefaultEdgePort = 6061
)

type Star struct {
	*option.StarConfig
	tap       *tuntap.Tuntap
	socket    socket.Interface
	cache     node.NodesCache //获取回来的Peers  mac: Node
	executor  map[int]executor.Executor
	processor sync.Map
	inbound   []chan *packet.Packet
}

func (star *Star) Start() error {
	once.Do(func() {
		star.socket = socket.NewSocket()
		if err := star.conn(); err != nil {
			logger.Errorf("failed to connect to registry: %v", err)
		}
		star.cache = node.New()
		star.Protocol = option.UDP
		tap, err := tuntap.New(tuntap.TAP)
		star.tap = tap

		if err != nil {
			logger.Errorf("create or connect tap failed, err: (%v)", err)
		}

		if err := star.register(); err != nil {
			logger.Errorf("registry failed. (%v)", err)
		}

		//star.initExecutor()
		star.initProcessor()
		go func() {
			for {
				star.queryPeer()
				//连通
				star.dialNode()
				time.Sleep(30 * time.Second)
			}
		}()
	})
	star.starLoop()
	return nil
}

func (star *Star) initProcessor() {
	deviceHandler := middleware.WithMiddlewares(device.New(), auth.Middleware())
	deviceProcessor := processordevice.New(star.tap, deviceHandler)
	udpHandler := middleware.WithMiddlewares(udp.New(), auth.Middleware())
	udpProcessor := processorudp.New(udpHandler)

	star.processor.Store(star.tap.Fd, deviceProcessor)
	star.processor.Store(star.socket.(socket.Socket).Fd, udpProcessor)
}
