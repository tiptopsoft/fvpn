package fvpnc

import (
	"fmt"
	"github.com/interstellar-cloud/star/pkg/addr"
	"github.com/interstellar-cloud/star/pkg/log"
	"github.com/interstellar-cloud/star/pkg/option"
	"github.com/interstellar-cloud/star/pkg/packet/peer"
	"github.com/interstellar-cloud/star/pkg/packet/register"
	"github.com/interstellar-cloud/star/pkg/processor"
	"github.com/interstellar-cloud/star/pkg/socket"
	"github.com/interstellar-cloud/star/pkg/util"
	"golang.org/x/sys/unix"
)

var logger = log.Log()

func (node *Node) conn() error {
	var err error
	switch node.Protocol {
	case option.UDP:
		remoteAddr, err := util.GetAddress(node.ClientCfg.Registry, addr.DefaultPort)
		if err != nil {
			return err
		}

		if err = node.socket.Connect(&remoteAddr); err != nil {
			return err
		}
		logger.Infof("node connected to fvpns: (%v)", node.ClientCfg.Registry)
	}
	return err
}

func (node *Node) queryPeer() error {
	cp := peer.NewPacket()
	data, err := cp.Encode()
	if err != nil {
		return err
	}

	switch node.Protocol {
	case option.UDP:
		logger.Infof("start to query node peer info, data: (%v)", data)
		if _, err := node.socket.Write(data); err != nil {
			return nil
		}
		break
	}
	return nil
}

// register register a edgestar to center.
func (node *Node) register() error {
	var err error
	rp := register.NewPacket()
	rp.SrcMac, _ = addr.GetMacAddrByDev(node.tap.Name)
	logger.Infof("register src mac: %v to fvpns", rp.SrcMac.String())
	data, err := rp.Encode()
	logger.Infof("sending fvpns data: %v", data)
	if err != nil {
		return err
	}
	switch node.Protocol {
	case option.UDP:
		logger.Infof("node start to register to fvpns: %v", rp)
		if _, err := node.socket.Write(data); err != nil {
			return err
		}
		break
	}
	return nil
}

// register register a edgestar to center.
func (node *Node) unregister() error {
	var err error
	rp := register.NewUnregisterPacket()
	rp.SrcMac = node.tap.MacAddr
	data, err := rp.Encode()
	fmt.Println("sending unregister data: ", data)
	if err != nil {
		return err
	}

	switch node.Protocol {
	case option.UDP:
		logger.Infof("node start to fvpns self to fvpns: %v", rp)
		if _, err := node.socket.Write(data); err != nil {
			return err
		}
		break
	}
	return nil
}

func (node *Node) starLoop() {
	netFd := node.socket.(socket.Socket).Fd
	tapFd := node.tap.Fd
	var FdSet unix.FdSet
	var maxFd int
	if netFd > tapFd {
		maxFd = netFd
	} else {
		maxFd = tapFd
	}
	for {
		FdSet.Zero()
		FdSet.Set(netFd)
		FdSet.Set(tapFd)
		timeout := &unix.Timeval{
			Sec:  3,
			Usec: 0,
		}
		ret, err := unix.Select(maxFd+1, &FdSet, nil, nil, timeout)
		if ret < 0 && err == unix.EINTR {
			continue
		}

		if err != nil {
			fmt.Println(err)
		}

		if FdSet.IsSet(tapFd) {
			if p, ok := node.processor.Load(tapFd); ok {
				if err := p.(processor.Processor).Process(); err != nil {
					logger.Errorf("tap process failed. %v", err)
				}
			} else {
				logger.Errorf("can not found tap socket")
			}
		}

		if FdSet.IsSet(netFd) {
			if p, ok := node.processor.Load(netFd); ok {
				if err := p.(processor.Processor).Process(); err != nil {
					logger.Errorf("net process failed. %v", err)
				}
			} else {
				logger.Errorf("can not found net socket")
			}
		}
	}
}

// use host socket write to destination, superNode or use p2p
func write2Net(socket socket.Interface, b []byte) {
	logger.Debugf("tap write to net packet: (%v)", b)
	if _, err := socket.Write(b); err != nil {
		logger.Errorf("tap write to net failed. (%v)", err)
	}
}

func (node Node) dialNode() {
	for _, v := range node.cache.Nodes {
		if v != nil && v.Addr != nil {
			dstAddr := v.Addr.(*unix.SockaddrInet4).Addr
			newAddr := &unix.SockaddrInet4{Addr: dstAddr, Port: DefaultEdgePort}
			if !v.P2P {
				if err := node.socket.Connect(newAddr); err != nil {
					return
				}
			}
			//如果连通，则更新cache中的状态
			v.P2P = true
		}

	}
}
