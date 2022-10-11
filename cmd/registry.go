package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/tiptopsoft/fvpn/pkg/relay"
	"github.com/tiptopsoft/fvpn/pkg/util"
)

type RegOptions struct {
	Listen int
}

func registryCmd() *cobra.Command {
	var opts RegOptions
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
	fs.IntVarP(&opts.Listen, "port", "p", 0, "registry server port")

	return cmd
}

func runSuper(opts *RegOptions) error {

	config, err := util.InitConfig()
	if err != nil {
		return err
	}
	s := relay.RegServer{
		RegistryCfg: config.RegistryCfg,
	}

	if opts.Listen != 0 {
		s.RegistryCfg.Listen = fmt.Sprintf(":%d", opts.Listen)
	}

	return s.Start()
}
