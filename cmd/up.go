package cmd

import (
	"github.com/interstellar-cloud/star/pkg/auth"
	"github.com/interstellar-cloud/star/pkg/edge"
	"github.com/spf13/cobra"
)

type upOptions struct {
	edge.EdgeConfig
	auth.StarAuth
	StarConfigFilePath string
}

func EdgeCmd() *cobra.Command {
	var opts upOptions
	var cmd = &cobra.Command{
		Use:          "up",
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
	_, err := edge.InitConfig()
	if err != nil {
		return err
	}

	s := edge.EdgeStar{}
	return s.Start()
}
