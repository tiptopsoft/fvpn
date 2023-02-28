package edge

import (
	"fmt"
	"github.com/interstellar-cloud/star/pkg/addr"
	"github.com/interstellar-cloud/star/pkg/log"
	"github.com/interstellar-cloud/star/pkg/option"
	"github.com/interstellar-cloud/star/pkg/packet/peer"
	"github.com/interstellar-cloud/star/pkg/packet/register"
	"github.com/interstellar-cloud/star/pkg/socket"
	"github.com/interstellar-cloud/star/pkg/util"
	"golang.org/x/sys/unix"
)

func (star *Star) conn() error {
	var err error
	switch star.Protocol {
	case option.UDP:
		remoteAddr, err := util.GetAddress(star.Registry, addr.DefaultPort)
		if err != nil {
			return err
		}

		if err = star.Socket.Connect(&remoteAddr); err != nil {
			return err
		}
		log.Logger.Infof("star connected to registry: (%v)", star.Registry)
	}
	return err
}

func (star *Star) queryPeer() error {
	cp := peer.NewPacket()
	data, err := cp.Encode()
	if err != nil {
		return err
	}

	switch star.Protocol {
	case option.UDP:
		log.Logger.Infof("start to query star peer info, data: (%v)", data)
		if _, err := star.Socket.Write(data); err != nil {
			return nil
		}
		break
	}
	return nil
}

// register register a edgestar to center.
func (star *Star) register() error {
	var err error
	rp := register.NewPacket()
	rp.SrcMac, _ = addr.GetMacAddrByDev(star.tap.Name)
	log.Logger.Infof("register src mac: %v to registry", rp.SrcMac.String())
	data, err := rp.Encode()
	log.Logger.Infof("sending registry data: %v", data)
	if err != nil {
		return err
	}
	switch star.Protocol {
	case option.UDP:
		log.Logger.Infof("star start to register to registry: %v", rp)
		if _, err := star.Socket.Write(data); err != nil {
			return err
		}
		break
	}
	return nil
}

// register register a edgestar to center.
func (star *Star) unregister() error {
	var err error

	rp := register.NewUnregisterPacket()
	rp.SrcMac = star.tap.MacAddr
	data, err := rp.Encode()
	fmt.Println("sending unregister data: ", data)
	if err != nil {
		return err
	}

	switch star.Protocol {
	case option.UDP:
		log.Logger.Infof("star start to registry self to registry: %v", rp)
		if _, err := star.Socket.Write(data); err != nil {
			return err
		}
		break
	}
	return nil
}

func (star *Star) starLoop() {
	netFd := star.Socket.(socket.Socket).Fd
	tapFd := star.tap.Fd
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
			if err := star.tapFunc(star.tap, star.Socket); err != nil {
				log.Logger.Errorf("tap socket failed. (%v)", err)
			}
		}

		if FdSet.IsSet(netFd) {
			if err := star.socketFunc(star.tap, star.Socket); err != nil {
				log.Logger.Errorf("socket func failed. (%v)", err)
			}
		}
	}
}

//use host socket write to destination, superNode or use p2p
func write2Net(socket socket.Interface, b []byte) {
	log.Logger.Debugf("tap write to net packet: (%v)", b)
	if _, err := socket.Write(b); err != nil {
		log.Logger.Errorf("tap write to net failed. (%v)", err)
	}
}

func (star Star) dialNode() {
	for _, v := range star.cache.Nodes {
		dstAddr := v.Addr.(*unix.SockaddrInet4).Addr
		newAddr := &unix.SockaddrInet4{Addr: dstAddr, Port: DefaultEdgePort}
		if !v.P2P {
			if err := star.Socket.Connect(newAddr); err != nil {
				return
			}
		}
		//如果连通，则更新cache中的状态
		v.P2P = true
	}
}
