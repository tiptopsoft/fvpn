package client

import (
	"github.com/topcloudz/fvpn/pkg/handler"
	"github.com/topcloudz/fvpn/pkg/handler/device"
	"github.com/topcloudz/fvpn/pkg/handler/udp"
	"github.com/topcloudz/fvpn/pkg/middleware/infra"
	"sync"

	"github.com/topcloudz/fvpn/pkg/middleware"
	"github.com/topcloudz/fvpn/pkg/option"
	"github.com/topcloudz/fvpn/pkg/socket"
	_ "github.com/topcloudz/fvpn/pkg/util"
)

var (
	once        sync.Once
	DefaultPort = 6663
)

type Node struct {
	*option.Config
	Protocol    option.Protocol
	tun         *handler.Tun //key: networkId, value: Tuntap
	relaySocket socket.Interface
}

func (n *Node) Start() error {
	once.Do(func() {
		n.relaySocket = socket.NewSocket(6061)
		n.Protocol = option.UDP
		if err := n.conn(); err != nil {
			logger.Errorf("failed to connect to server: %v", err)
		}
	})
	tun := n.GetTun() //这里启动的是relaySocket，中继服务器
	go tun.ReadFromUdp()

	//go tun.QueryRemoteNodes()
	go tun.WriteToDevice()
	//open hole for p2p
	go tun.PunchHole()

	return n.runHttpServer()

}

func (n *Node) GetTun() *handler.Tun {
	m := n.initMiddleware()
	tunHandler := middleware.WithMiddlewares(device.Handle(), m...)
	udpHandler := middleware.WithMiddlewares(udp.Handle(), m...)
	tun := handler.NewTun(tunHandler, udpHandler, n.relaySocket)
	n.tun = tun
	return tun
}

// initMiddleware TODO add impl
func (n *Node) initMiddleware() []middleware.Middleware {
	return infra.Middlewares(n.OpenAuth, n.OpenEncrypt)
}
