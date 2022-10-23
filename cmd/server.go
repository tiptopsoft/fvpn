package cmd

import (
	"github.com/interstellar-cloud/star/common"
	"github.com/interstellar-cloud/star/option"
	"github.com/interstellar-cloud/star/user"
	"github.com/spf13/cobra"
)

type serverOptions struct {
	option.StarConfig
	*option.Config
}

func serverCmd() *cobra.Command {
	var opts serverOptions
	cmd := &cobra.Command{
		Use:   "serve",
		Short: "remove a device",
		PreRunE: func(cmd *cobra.Command, args []string) error {

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			return runServer(&opts)
		},
	}

	fs := cmd.Flags()
	fs.IntVarP(&opts.Port, "port", "p", 3000, "tun server port")

	return cmd
}

func runServer(opts *serverOptions) error {
	config, err := common.InitConfig()
	if err != nil {
		return err
	}
	opts.Config = config

	s := user.Server{
		Config: config,
	}

	return s.Start()
}
