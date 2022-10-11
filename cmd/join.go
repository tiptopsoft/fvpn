package cmd

import (
	"errors"
	"fmt"
	"github.com/spf13/cobra"
	"github.com/tiptopsoft/fvpn/pkg/device"
	"github.com/tiptopsoft/fvpn/pkg/util"
)

type joinOptions struct {
	*util.NodeCfg
	StarConfigFilePath string
	addr               string
	networkId          string
}

func joinCmd() *cobra.Command {
	var opts joinOptions
	var cmd = &cobra.Command{
		Use:          "join",
		SilenceUsage: true,
		Short:        "join a network",
		Long: `join a network which created by user, 
networkId could be found on our site after user registered, 
use free services or pay services`,

		PreRunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},

		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("networkId should be given")
			}
			return runJoin(args, &opts)
		},
	}
	fs := cmd.Flags()
	fs.StringVarP(&opts.networkId, "id", "id", "", "private network id")

	return cmd
}

// runJoin join a network cmd
func runJoin(args []string, opts *joinOptions) error {
	cfg, err := util.InitConfig()
	if err != nil {
		return err
	}

	var networkId string
	if args[0] != "" {
		networkId = args[0]
	} else {
		networkId = opts.networkId
	}

	if networkId == "" {
		return errors.New("networkId is empty")
	}
	if err := device.RunJoinNetwork(cfg, networkId); err != nil {
		return err
	}

	fmt.Println(fmt.Sprintf("Join to network: %s successed", networkId))
	return nil
}
