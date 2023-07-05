package cmd

import (
	"github.com/spf13/cobra"
)

type RegStarOptions struct {
	Listen int
}

func RegCmd() *cobra.Command {
	var opts RegStarOptions
	cmd := &cobra.Command{
		Use:   "serve",
		Short: "s",
		Long:  `fvpn start a registry, a relay server for node to node mesh.`,
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

	//config, err := util.InitConfig()
	//if err != nil {
	//	return err
	//}
	//s := server.RegServer{
	//	ServerConfig: config.ServerCfg,
	//	Outbound:     make(chan *packet.Frame, 15000),
	//}

	return nil
}
