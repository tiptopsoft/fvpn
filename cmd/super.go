package cmd

import (
	"github.com/interstellar-cloud/star/pkg/option"
	"github.com/interstellar-cloud/star/pkg/super"
	"github.com/spf13/cobra"
)

type superStarOptions struct {
	Config *option.Config
	Listen int
}

func superStarCmd() *cobra.Command {
	var opts superStarOptions
	cmd := &cobra.Command{
		Use: "super",
		Short: `super star, using for finding other machine in a group,
which star can register to, also super can super packets when star at a Symetric Nat.`,
		PreRunE: func(cmd *cobra.Command, args []string) error {

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			return runSuper(&opts)
		},
	}

	fs := cmd.Flags()
	fs.IntVarP(&opts.Listen, "port", "p", 3000, "tun server port")

	return cmd
}

func runSuper(opts *superStarOptions) error {
	config, err := option.InitConfig()
	if err != nil {
		return err
	}
	opts.Config = config

	s := super.RelayServer{
		Config: config,
	}

	return s.Start(opts.Listen)
}
