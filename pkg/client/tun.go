package client

import (
	"github.com/topcloudz/fvpn/pkg/handler/device"
	"github.com/topcloudz/fvpn/pkg/middleware"
	processordevice "github.com/topcloudz/fvpn/pkg/processor/device"
	"github.com/topcloudz/fvpn/pkg/tuntap"
)

var (
	limit = 10
)

// 通过netId获取tuntap
func (n *Node) Join(netId string) error {
	tun, err := tuntap.GetTuntap(netId)
	if err != nil {
		return err
	}
	n.taps.Store(netId, tun)
	n.newProcessor(tun)
	//注册tun到registry
	return n.register(tun)
}

func (n *Node) newProcessor(tun *tuntap.Tuntap) {
	deviceHandler := middleware.WithMiddlewares(device.New(tun, n.socket, n.cache), n.initMiddleware()...)
	deviceProcessor := processordevice.New(tun, deviceHandler)
	n.processor.Store(tun.Fd, deviceProcessor)
}
