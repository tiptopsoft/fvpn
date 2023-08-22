package cmd

import (
	"github.com/spf13/cobra"
	"github.com/tiptopsoft/fvpn/pkg/device"
	"github.com/tiptopsoft/fvpn/pkg/util"
	"log"
)

type upOptions struct {
	*util.NodeCfg
	Daemon bool
}

func nodeCmd() *cobra.Command {
	var opts upOptions
	var cmd = &cobra.Command{
		Use:          "node",
		Aliases:      []string{"n"},
		SilenceUsage: true,
		Short:        "start up a node, for private network proxy",
		Long:         `start up a node is start a private network proxy, use fvpn, you can use any device visit your private network from any place`,

		PreRunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},

		RunE: func(cmd *cobra.Command, args []string) error {
			if opts.Daemon {
				cntxt := util.GetDaemon()

				d, err := cntxt.Reborn()
				if err != nil {
					log.Fatal("Unable to run: ", err)
				}
				if d != nil {
					return nil
				}
				defer cntxt.Release()

				log.Print("fvpn started")
			}
			return runNode(&opts)
		},
	}

	fs := cmd.Flags()
	fs.BoolVarP(&opts.Daemon, "daemon", "d", false, "run daemon")

	return cmd
}

// runEdge run a client up
func runNode(opts *upOptions) error {
	config, err := util.InitConfig()
	if err != nil {
		return err
	}

	return device.Start(config)
}
