package cmd

import (
	"github.com/interstellar-cloud/star/pkg/option"
	"github.com/interstellar-cloud/star/pkg/user"
	"github.com/spf13/cobra"
)

type serverOptions struct {
	*option.Config
}

func SuperCmd() *cobra.Command {
	var opts serverOptions
	cmd := &cobra.Command{
		Use:   "serve",
		Short: "remove a tuntap",
		PreRunE: func(cmd *cobra.Command, args []string) error {

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			return runServer(&opts)
		},
	}

	//fs := cmd.Flags()
	//fs.IntVarP(&opts.Port, "port", "p", 3000, "tun server port")

	return cmd
}

func runServer(opts *serverOptions) error {
	config, err := option.InitConfig()
	if err != nil {
		return err
	}
	opts.Config = config

	s := user.UserServer{
		Config: config,
	}

	return s.Start(3000)
}
