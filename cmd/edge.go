package cmd

import (
	"github.com/interstellar-cloud/star/pkg/edge"
	"github.com/interstellar-cloud/star/pkg/handler/auth"
	option2 "github.com/interstellar-cloud/star/pkg/util/option"
	"github.com/spf13/cobra"
)

type upOptions struct {
	*option2.EdgeConfig
	auth.StarAuth
	StarConfigFilePath string
}

func EdgeCmd() *cobra.Command {
	var opts upOptions
	var cmd = &cobra.Command{
		Use:          "edge",
		SilenceUsage: true,
		Short:        "start up a edge, for net proxy",
		Long:         `Start up a edge, for private net proxy`,

		PreRunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},

		RunE: func(cmd *cobra.Command, args []string) error {
			return runEdge(&opts)
		},
	}
	fs := cmd.Flags()
	fs.StringVarP(&opts.StarConfigFilePath, "config", "", "", "config file for edge")

	return cmd
}

//runEdge run a edge up
func runEdge(opts *upOptions) error {
	config, err := option2.InitConfig()
	if err != nil {
		return err
	}

	s := edge.StarEdge{
		EdgeConfig: config.Star,
	}
	return s.Start()
}
