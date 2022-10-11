package cmd

import (
	"github.com/spf13/cobra"
	"github.com/tiptopsoft/fvpn/pkg/device"
	"github.com/tiptopsoft/fvpn/pkg/util"
)

type upOptions struct {
	*util.NodeCfg
	Daemon bool
}

func nodeCmd() *cobra.Command {
	var opts upOptions
	var cmd = &cobra.Command{
		Use:          "node",
		Aliases:      []string{"n"},
		SilenceUsage: true,
		Short:        "start up a node, for private network proxy",
		Long:         `start up a node is start a private network proxy, use fvpn, you can use any device visit your private network from any place`,

		PreRunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},

		RunE: func(cmd *cobra.Command, args []string) error {
			return runNode(&opts)
		},
	}

	return cmd
}

// runEdge run a client up
func runNode(opts *upOptions) error {
	config, err := util.InitConfig()
	if err != nil {
		return err
	}

	return device.Start(config)
}
