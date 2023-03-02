package edge

import (
	"github.com/interstellar-cloud/star/pkg/executor"
	"github.com/interstellar-cloud/star/pkg/node"
	"github.com/interstellar-cloud/star/pkg/option"
	"github.com/interstellar-cloud/star/pkg/packet"
	"github.com/interstellar-cloud/star/pkg/socket"
	"github.com/interstellar-cloud/star/pkg/tuntap"
	"sync"
	"time"
)

var (
	once            sync.Once
	DefaultEdgePort = 6061
)

type Star struct {
	*option.StarConfig
	tap      *tuntap.Tuntap
	socket   socket.Interface
	cache    node.NodesCache //获取回来的Peers  mac: Node
	executor map[int]executor.Executor
	inbound  []chan *packet.Packet
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

		star.initExecutor()
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

func (star *Star) initExecutor() {
	star.executor = make(map[int]executor.Executor, 1)
	t := TapExecutor{
		device: star.tap,
		cache:  star.cache,
	}

	star.executor[star.tap.Fd] = t

	s := SocketExecutor{
		device:   star.tap,
		Protocol: star.Protocol,
		cache:    star.cache,
	}
	star.executor[star.socket.(socket.Socket).Fd] = s
}
