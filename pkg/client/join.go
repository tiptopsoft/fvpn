package client

import (
	"errors"
	"github.com/topcloudz/fvpn/pkg/http"
	"github.com/topcloudz/fvpn/pkg/util"
)

const (
	userUrl  = "http://211.159.225.186"
	localUrl = "http://127.0.0.1:6663"
)

func (p *Peer) RunJoinNetwork(netId string) error {
	logger.Infof("start to join %s", netId)
	req := new(http.JoinRequest)
	mac, err := util.GetHostMac()
	if err != nil {
		return errors.New("can not found default host mac")
	}
	req.SrcMac = mac.String()
	req.NetworkId = netId

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

func (p *Peer) RunLeaveNetwork(networkId string) error {

	return nil
}
