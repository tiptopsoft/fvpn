package client

import (
	"fmt"
	"github.com/topcloudz/fvpn/pkg/addr"
	"github.com/topcloudz/fvpn/pkg/log"
	"github.com/topcloudz/fvpn/pkg/option"
	"github.com/topcloudz/fvpn/pkg/packet/register"
	"github.com/topcloudz/fvpn/pkg/tuntap"
	"github.com/topcloudz/fvpn/pkg/util"
)

var (
	logger = log.Log()
)

func (n *Node) conn() error {
	var err error
	switch n.Protocol {
	case option.UDP:
		remoteAddr, err := util.GetAddress(n.ClientCfg.Registry, addr.DefaultPort)
		if err != nil {
			return err
		}

		if err = n.relaySocket.Connect(&remoteAddr); err != nil {
			return err
		}

		n.relayAddr = &remoteAddr
		logger.Infof("node connected to server: (%v)", n.ClientCfg.Registry)
	}
	return err
}

// register register a edgestar to center.
func (n *Node) unregister(tun *tuntap.Tuntap) error {
	var err error
	rp := register.NewUnregisterPacket(tun.NetworkId)
	copy(rp.SrcMac[:], tun.MacAddr)
	data, err := register.Encode(rp)
	fmt.Println("sending unregister data: ", data)
	if err != nil {
		return err
	}

	switch n.Protocol {
	case option.UDP:
		logger.Infof("node start to server self to server: %v", rp)
		if _, err := n.relaySocket.Write(data); err != nil {
			return err
		}
		break
	}
	return nil
}
