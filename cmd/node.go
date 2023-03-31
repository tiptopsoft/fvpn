package cmd

import (
	"github.com/interstellar-cloud/star/pkg/fvpnc"
	"github.com/interstellar-cloud/star/pkg/middleware/auth"
	option2 "github.com/interstellar-cloud/star/pkg/option"
	"github.com/spf13/cobra"
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

// runEdge run a fvpnc up
func runNode(opts *upOptions) error {
	config, err := option2.InitConfig()
	if err != nil {
		return err
	}

	s := &fvpnc.Node{
		Config: config,
	}
	return s.Start()
}
