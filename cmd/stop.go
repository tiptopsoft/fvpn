package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/tiptopsoft/fvpn/pkg/device"
	"github.com/tiptopsoft/fvpn/pkg/util"
)

func stopCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "stop",
		Short: "stop a running node",
		Long:  "stop a running node",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runStop()
		},
	}

	return cmd
}

func runStop() error {
	cfg, err := util.InitConfig()
	if err != nil {
		return err
	}
	err = device.Stop(cfg.NodeCfg)
	if err != nil {
		return err
	}

	fmt.Println("closed fvpn success.")
	return nil
}
