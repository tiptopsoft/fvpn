package client

import (
	"fmt"
	"github.com/topcloudz/fvpn/pkg/addr"
	"github.com/topcloudz/fvpn/pkg/log"
	"github.com/topcloudz/fvpn/pkg/option"
	"github.com/topcloudz/fvpn/pkg/socket"
)

var (
	logger = log.Log()
)

func (p *Peer) conn() error {
	var err error
	switch p.Protocol {
	case option.UDP:
		if s, err := socket.NewSocket("", fmt.Sprintf("%s:%d", p.ClientCfg.Registry, addr.DefaultPort)); err != nil {
			return err
		} else {
			p.relaySocket = s
		}
		logger.Infof("node connected to server: (%v)", p.ClientCfg.Registry)
	}
	return err
}
