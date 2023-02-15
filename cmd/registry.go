package cmd

import (
	"github.com/interstellar-cloud/star/pkg/registry"
	"github.com/interstellar-cloud/star/pkg/util/handler/auth"
	"github.com/interstellar-cloud/star/pkg/util/handler/encrypt"
	"github.com/interstellar-cloud/star/pkg/util/option"
	"github.com/spf13/cobra"
)

type RegStarOptions struct {
	Listen int
}

func RegCmd() *cobra.Command {
	var opts RegStarOptions
	cmd := &cobra.Command{
		Use:   "registry",
		Short: "reg",
		Long: `registry edge, using for finding other machine in a group,
which edge can registry to, also registry can registry packets when edge at a Symetric Nat.`,
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

func runSuper(opts *RegStarOptions) error {

	config, err := option.InitConfig()
	if err != nil {
		return err
	}
	s := registry.RegStar{
		RegConfig:      config.Reg,
		AuthHandler:    new(auth.AuthHandler),
		EncryptHandler: new(encrypt.EncryptHandler),
	}

	return s.Start(config.Reg.Listen)
}
