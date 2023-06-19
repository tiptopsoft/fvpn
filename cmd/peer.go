package cmd

import (
	"github.com/spf13/cobra"
	"github.com/topcloudz/fvpn/pkg/client"
	"github.com/topcloudz/fvpn/pkg/middleware/auth"
	option2 "github.com/topcloudz/fvpn/pkg/option"
)

type upOptions struct {
	*option2.ClientConfig
	auth.StarAuth
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
	config, err := option2.InitConfig()
	if err != nil {
		return err
	}

	s := &client.Peer{
		Config: config,
	}
	return s.Start()
}
