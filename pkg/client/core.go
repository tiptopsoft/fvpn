package client

import (
	"github.com/topcloudz/fvpn/pkg/addr"
	"github.com/topcloudz/fvpn/pkg/log"
	"github.com/topcloudz/fvpn/pkg/option"
	"github.com/topcloudz/fvpn/pkg/util"
)

var (
	logger = log.Log()
)

func (p *Peer) conn() error {
	var err error
	switch p.Protocol {
	case option.UDP:
		remoteAddr, err := util.GetAddress(p.ClientCfg.Registry, addr.DefaultPort)
		if err != nil {
			return err
		}

		if err = p.relaySocket.Connect(&remoteAddr); err != nil {
			return err
		}

		p.relayAddr = &remoteAddr
		logger.Infof("node connected to server: (%v)", p.ClientCfg.Registry)
	}
	return err
}
