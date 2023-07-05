package cmd

import (
	"github.com/spf13/cobra"
	"github.com/topcloudz/fvpn/pkg/node"
	"github.com/topcloudz/fvpn/pkg/util"
)

type upOptions struct {
	*util.ClientConfig
	StarConfigFilePath string
}

func EdgeCmd() *cobra.Command {
	var opts upOptions
	var cmd = &cobra.Command{
		Use:          "node",
		Aliases:      []string{"n"},
		SilenceUsage: true,
		Short:        "start up a node, for private net proxy",
		Long:         `Start up a node, for private net proxy`,

		PreRunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},

		RunE: func(cmd *cobra.Command, args []string) error {
			return runNode(&opts)
		},
	}
	fs := cmd.Flags()
	fs.StringVarP(&opts.StarConfigFilePath, "config", "", "", "config file for fvpn")

	return cmd
}

// runEdge run a client up
func runNode(opts *upOptions) error {
	config, err := util.InitConfig()
	if err != nil {
		return err
	}

	return node.Start(config)
}
