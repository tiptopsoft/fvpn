package node

import (
	"github.com/topcloudz/fvpn/pkg/util"
)

func RunJoinNetwork(cfg *util.Config, networkId string) error {
	logger.Infof("start join to network: %s", networkId)

	cm := NewManager(cfg.ClientCfg)
	resp, err := cm.JoinNetwork(networkId)
	if err != nil {
		return err
	}

	return NewRouter(resp.IP, resp.Name).AddRouter(resp.Network)
}

func (p *Peer) RunLeaveNetwork(networkId string) error {

	return nil
}
