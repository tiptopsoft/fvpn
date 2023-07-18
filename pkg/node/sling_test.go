package node

import (
	"github.com/topcloudz/fvpn/pkg/util"
	"testing"
)

func TestClient_Init(t *testing.T) {
	cfg, _ := util.InitConfig()
	client := NewClient(cfg.ClientCfg.ControlUrl())
	client.Init()
}
