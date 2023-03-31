package cmd

import (
	"github.com/interstellar-cloud/star/pkg/option"
	"github.com/interstellar-cloud/star/pkg/registry"
	"github.com/spf13/cobra"
)

type regStarOptions struct {
	Listen int
}

func regCmd() *cobra.Command {
	var opts regStarOptions
	cmd := &cobra.Command{
		Use:   "registry",
		Short: "reg",
		Long: `registry client, using for finding other machine in a group,
which client can registry to, also registry can registry packets when client at a Symetric Nat.`,
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

func runSuper(opts *regStarOptions) error {

	config, err := option.InitConfig()
	if err != nil {
		return err
	}
	s := registry.RegStar{
		RegConfig: config.Reg,
	}

	return s.Start(config.Reg.Listen)
}
