package client

import (
	"github.com/topcloudz/fvpn/pkg/handler"
	"github.com/topcloudz/fvpn/pkg/handler/device"
	"github.com/topcloudz/fvpn/pkg/handler/udp"
	"github.com/topcloudz/fvpn/pkg/middleware/encrypt"
	"sync"

	"github.com/topcloudz/fvpn/pkg/middleware"
	"github.com/topcloudz/fvpn/pkg/middleware/auth"
	"github.com/topcloudz/fvpn/pkg/option"
	"github.com/topcloudz/fvpn/pkg/socket"
)

var (
	once        sync.Once
	DefaultPort = 6663
)

type Node struct {
	*option.Config
	Protocol    option.Protocol
	tuns        sync.Map //key: netId, value: Tuntap
	relaySocket socket.Interface
}

func (n *Node) Start() error {
	once.Do(func() {
		n.relaySocket = socket.NewSocket()
		if err := n.conn(); err != nil {
			logger.Errorf("failed to connect to server: %v", err)
		}
		n.Protocol = option.UDP
	})
	tun := n.GetTun() //这里启动的是relaySocket，中继服务器
	go tun.ReadFromUdp()
	go tun.WriteToDevice()
	return n.runHttpServer()
}

func (n *Node) GetTun() handler.Tun {
	m := n.initMiddleware()
	tunHandler := middleware.WithMiddlewares(device.Handle(), m...)
	udpHandler := middleware.WithMiddlewares(udp.Handle(), m...)
	tun := handler.NewTun(tunHandler, udpHandler, n.relaySocket, 0)
	return *tun
}

// initMiddleware TODO add impl
func (n *Node) initMiddleware() []middleware.Middleware {
	var result []middleware.Middleware
	if n.OpenAuth {
		result = append(result, auth.Middleware())
	}

	if n.OpenEncrypt {
		result = append(result, encrypt.Middleware())
	}

	if n.OpenCompress {
		//TODO
	}

	return result
}
