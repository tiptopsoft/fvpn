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

		if err = n.socket.Connect(&remoteAddr); err != nil {
			return err
		}
		logger.Infof("n connected to server: (%v)", n.ClientCfg.Registry)
	}
	return err
}

func (n *Node) queryPeer() error {
	cp := peer.NewPacket()
	data, err := cp.Encode()
	if err != nil {
		return err
	}

	switch n.Protocol {
	case option.UDP:
		logger.Infof("start to query n peer info, data: (%v)", data)
		if _, err := n.socket.Write(data); err != nil {
			return nil
		}
		break
	}
	return nil
}

// register register a edgestar to center.
func (n *Node) register(tun *tuntap.Tuntap) error {
	var err error
	rp := register.NewPacket()
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
		if _, err := n.socket.Write(data); err != nil {
			return err
		}
		break
	}
	return nil
}

// register register a edgestar to center.
func (n *Node) unregister(tun *tuntap.Tuntap) error {
	var err error
	rp := register.NewUnregisterPacket()
	rp.SrcMac = tun.MacAddr
	data, err := rp.Encode()
	fmt.Println("sending unregister data: ", data)
	if err != nil {
		return err
	}

	switch n.Protocol {
	case option.UDP:
		logger.Infof("n start to server self to server: %v", rp)
		if _, err := n.socket.Write(data); err != nil {
			return err
		}
		break
	}
	return nil
}

//func (n *Node) dialNode() {
//	for _, v := range n.cache.Nodes {
//		if v != nil && v.Addr != nil {
//			dstAddr := v.Addr.(*unix.SockaddrInet4).Addr
//			newAddr := &unix.SockaddrInet4{Addr: dstAddr, Port: DefaultPort}
//			if !v.P2P {
//				if err := n.socket.Connect(newAddr); err != nil {
//					return
//				}
//			}
//			//如果连通，则更新cache中的状态
//			v.P2P = true
//		}
//	}
//}
