package cmd

import (
	"github.com/spf13/cobra"
	"github.com/tiptopsoft/fvpn/pkg/device"
	"github.com/tiptopsoft/fvpn/pkg/util"
)

func statusCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "status",
		Short: "check fvpn status",
		Long:  "check fvpn status, if fvpn node is running will return a message, if not, will return other.",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runStatus()
		},
	}

	return cmd
}

func runStatus() error {
	cfg, err := util.InitConfig()
	if err != nil {
		return err
	}
	return device.Status(cfg.NodeCfg)
}
