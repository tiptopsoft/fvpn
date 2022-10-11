package cmd

import (
	"errors"
	"github.com/spf13/cobra"
	"github.com/tiptopsoft/fvpn/pkg/device"
	"github.com/tiptopsoft/fvpn/pkg/util"
)

// Join a networkId will be created tun device. and will be assigned a IP which can be found in our website.
type leaveOptions struct {
	*util.NodeCfg
	StarConfigFilePath string
	NetworkId          string
}

func leaveCmd() *cobra.Command {
	var opts joinOptions
	var cmd = &cobra.Command{
		Use:          "leave",
		SilenceUsage: true,
		Short:        "leave a network",
		Long: `leave a joined network, once use leave a network, 
fvpn can not route any frame to dst node, 
if you want continue your destination routing, 
you can join it again`,

		PreRunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},

		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("networkId should be given")
			}
			return runLeave(args)
		},
	}
	fs := cmd.Flags()
	fs.StringVarP(&opts.networkId, "id", "id", "", "network id")

	return cmd
}

// runJoin join a network cmd
func runLeave(args []string) error {

	cfg, err := util.InitConfig()
	if err != nil {
		return err
	}
	return device.RunLeaveNetwork(cfg, args[0])
}
