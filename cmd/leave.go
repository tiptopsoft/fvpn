package cmd

import (
	"errors"
	"github.com/spf13/cobra"
	"github.com/topcloudz/fvpn/pkg/client"
	"github.com/topcloudz/fvpn/pkg/middleware/auth"
	"github.com/topcloudz/fvpn/pkg/option"
)

// Join a networkId will be created tun device. and will be assigned a IP which can be found in our website.
type leaveOptions struct {
	*option.ClientConfig
	auth.StarAuth
	StarConfigFilePath string
	NetworkId          string
}

func leaveCmd() *cobra.Command {
	var opts joinOptions
	var cmd = &cobra.Command{
		Use:          "leave",
		SilenceUsage: true,
		Short:        "leave a network",
		Long:         `leave a network which created by user, network after created will show on our site`,

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
	fs.StringVarP(&opts.NetworkId, "config", "", "", "config file for fvpn")

	return cmd
}

// runJoin join a network cmd
func runLeave(args []string) error {
	config, err := option.InitConfig()
	if err != nil {
		return err
	}

	s := &client.Node{
		Config: config,
	}

	return s.RunJoinNetwork(args[0])
}
