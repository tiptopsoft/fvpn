package device

import (
	"fmt"
	"github.com/tiptopsoft/fvpn/pkg/util"
)

func Status(cfg *util.NodeCfg) error {
	client := NewClient(cfg.HostUrl())
	resp, err := client.Status()
	if err != nil {
		return err
	}

	if resp == nil || resp.Status == "" {
		fmt.Println("fvpn not running, please check")
	} else {
		fmt.Println(fmt.Sprintf("Status: %s, Version: %s", resp.Status, resp.Version))
	}
	return nil
}

func Stop(cfg *util.NodeCfg) error {
	client := NewClient(cfg.HostUrl())
	resp, err := client.Stop()
	if err != nil {
		return err
	}

	fmt.Println(resp.Result)
	return nil
}
