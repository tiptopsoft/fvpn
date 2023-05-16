package client

import (
	"errors"
	"github.com/topcloudz/fvpn/pkg/addr"
	"github.com/topcloudz/fvpn/pkg/http"
)

const (
	userUrl  = "https://www.efvpn.com"
	localUrl = "http://localhost:6663"
)

func (n *Node) RunJoinNetwork(netId string) error {
	logger.Infof("start to join %s", netId)
	req := new(http.JoinRequest)
	mac, err := addr.GetHostMac()
	if err != nil {
		return errors.New("can not found default host mac")
	}
	req.SrcMac = mac.String()
	req.NetworkId = netId
	req.UserId = "1"

	regClient := http.New(userUrl)
	resp, err := regClient.JoinNetwork(*req)
	if err != nil {
		return err
	}
	localClient := http.New(localUrl)
	req.NetworkId = resp.NetworkId
	req.Ip = resp.IP
	req.Mask = resp.Mask
	err = localClient.JoinLocalFvpn(*req)

	if err != nil {
		return err
	}

	logger.Infof("join network %s success", resp.NetworkId)

	return nil
}

func (n *Node) RunLeaveNetwork(networkId string) error {

	return nil
}
