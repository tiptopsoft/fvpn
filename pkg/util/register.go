package util

import (
	"errors"
	"github.com/topcloudz/fvpn/pkg/addr"
	"github.com/topcloudz/fvpn/pkg/packet/peer"
	"github.com/topcloudz/fvpn/pkg/packet/register"
	"github.com/topcloudz/fvpn/pkg/socket"
	"github.com/topcloudz/fvpn/pkg/tuntap"
	"log"
)

// register register a node to center.
func SendRegister(tun *tuntap.Tuntap, socket socket.Interface) error {
	var err error
	//header, err := packet.NewHeader(option.MsgTypeRegisterSuper, tun.NetworkId)
	srcMac, srcIP, err := addr.GetMacAddrAndIPByDev(tun.Name)
	if err != nil {
		return err
	}

	if srcIP == nil {
		return errors.New("device ip not set")
	}
	regPkt := register.NewPacket(tun.NetworkId, srcMac, srcIP)
	copy(regPkt.SrcMac[:], tun.MacAddr)
	if err != nil {
		return err
	}

	data, err := register.Encode(regPkt)
	if err != nil {
		return err
	}

	if _, err := socket.Write(data); err != nil {
		return err
	}
	logger.Debugf("sending server register data: %v", data)
	return nil
}

func SendQueryPeer(networkId string, socket socket.Interface) error {
	pkt := peer.NewPacket(networkId)
	buff, err := peer.Encode(pkt)
	if err != nil {
		log.Printf("query data failed: %v", err)
	}

	_, err = socket.Write(buff)
	if err != nil {
		return err
	}
	logger.Debugf("sending server query nodes data: %v", buff)

	return nil
}
