package cmd

import (
	"github.com/spf13/cobra"
	"github.com/tiptopsoft/fvpn/pkg/relay"
	"github.com/tiptopsoft/fvpn/pkg/util"
)

type RegStarOptions struct {
	Listen int
}

func RegCmd() *cobra.Command {
	var opts RegStarOptions
	cmd := &cobra.Command{
		Use:   "serve",
		Short: "s",
		Long:  `fvpn start a registry, a data center/relay server, is our core service`,
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

	config, err := util.InitConfig()
	if err != nil {
		return err
	}
	s := relay.RegServer{
		RegistryCfg: config.RegistryCfg,
	}

	return s.Start()
}
