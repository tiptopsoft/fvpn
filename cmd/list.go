package cmd

import (
	"github.com/spf13/cobra"
	"github.com/tiptopsoft/fvpn/pkg/device"
	"github.com/tiptopsoft/fvpn/pkg/util"
)

type listOptions struct {
}

func listCmd() *cobra.Command {
	opts := listOptions{}
	cmd := &cobra.Command{
		Use:          "list",
		SilenceUsage: true,
		Short:        "list",
		Long:         `when you've login in, list will show your networkIds in pretty json'`,

		PreRunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},

		RunE: func(cmd *cobra.Command, args []string) error {
			return runList(opts)
		},
	}

	return cmd
}

func runList(options listOptions) error {
	cfg, err := util.InitConfig()
	if err != nil {
		return err
	}
	return device.RunListNetworks(cfg)
}
