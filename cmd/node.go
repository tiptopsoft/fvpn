package cmd

import (
	"github.com/interstellar-cloud/star/pkg/client"
	"github.com/interstellar-cloud/star/pkg/middleware/auth"
	option2 "github.com/interstellar-cloud/star/pkg/option"
	"github.com/spf13/cobra"
)

type upOptions struct {
	*option2.StarConfig
	auth.StarAuth
	StarConfigFilePath string
}

func nodeCmd() *cobra.Command {
	var opts upOptions
	var cmd = &cobra.Command{
		Use:          "cache",
		SilenceUsage: true,
		Short:        "start up a cache, for private net proxy",
		Long:         `Start up a cache, for private net proxy`,

		PreRunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},

		RunE: func(cmd *cobra.Command, args []string) error {
			return runEdge(&opts)
		},
	}
	fs := cmd.Flags()
	fs.StringVarP(&opts.StarConfigFilePath, "config", "", "", "config file for client")

	return cmd
}

// runEdge run a client up
func runEdge(opts *upOptions) error {
	config, err := option2.InitConfig()
	if err != nil {
		return err
	}

	s := &client.Node{
		StarConfig: config.Star,
	}
	return s.Start()
}
