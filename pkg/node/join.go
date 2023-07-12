package node

import (
	"github.com/topcloudz/fvpn/pkg/util"
	"strings"
)

const (
	userUrl  = "http://211.159.225.186"
	localUrl = "http://127.0.0.1:6663"
)

func RunJoinNetwork(addr string) error {
	logger.Infof("start to join %s", addr)
	client := NewClient(localUrl)
	request := util.JoinRequest{
		IP:      "",
		Network: "",
	}

	if strings.Contains(addr, "/") {
		request.Network = addr
	} else {
		request.IP = addr
	}
	resp, err := client.JoinLocalFvpn(request)
	if err != nil {
		return err
	}

	return NewRouter(resp.IP, resp.Name).AddRouter(addr)
}

func (p *Peer) RunLeaveNetwork(networkId string) error {

	return nil
}
