package client

import (
	"fmt"
	"github.com/topcloudz/fvpn/pkg/addr"
	"github.com/topcloudz/fvpn/pkg/log"
	"github.com/topcloudz/fvpn/pkg/option"
	"github.com/topcloudz/fvpn/pkg/packet/peer"
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
		logger.Infof("n connected to server: (%v)", n.ClientCfg.Registry)
	}
	return err
}

func (n *Node) queryPeer() error {
	n.tuns.Range(func(key, value any) bool {
		networkId := key
		cp := peer.NewPacket(networkId.(string))
		data, err := cp.Encode()
		if err != nil {
			logger.Errorf("error occurd when query peers, networkId: %s, err: %v", networkId, err)
			return false
		}

		switch n.Protocol {
		case option.UDP:
			logger.Infof("start to query n peer info, data: (%v)", data)
			if _, err := n.relaySocket.Write(data); err != nil {
				logger.Errorf("error occurd when query peers, networkId: %s, err: %v", networkId, err)
				return false
			}
			break
		}
		return true
	})

	return nil
}

// register register a edgestar to center.
func (n *Node) register(tun *tuntap.Tuntap) error {
	var err error
	rp := register.NewPacket(tun.NetworkId)
	rp.SrcMac, _ = addr.GetMacAddrByDev(tun.Name)
	logger.Infof("register src mac: %v to server", rp.SrcMac.String())
	data, err := rp.Encode()
	logger.Infof("sending server data: %v", data)
	if err != nil {
		return err
	}
	switch n.Protocol {
	case option.UDP:
		logger.Infof("n start to register to server: %v", rp)
		if _, err := n.relaySocket.Write(data); err != nil {
			return err
		}
		break
	}
	return nil
}

// register register a edgestar to center.
func (n *Node) unregister(tun *tuntap.Tuntap) error {
	var err error
	rp := register.NewUnregisterPacket(tun.NetworkId)
	rp.SrcMac = tun.MacAddr
	data, err := rp.Encode()
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
