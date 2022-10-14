package cmd

import (
	"fmt"
	"github.com/interstellar-cloud/star/device"
	"github.com/interstellar-cloud/star/option"
	"github.com/spf13/cobra"
)

type upOptions struct {
	option.StarConfig
	option.StarAuth

	StarConfigFilePath string
}

func upCmd() *cobra.Command {
	var opts upOptions
	var cmd = &cobra.Command{
		Use:          "up",
		SilenceUsage: true,
		Short:        "start up a star, for net proxy",
		Long:         `Start up a star, for private net proxy`,

		PreRunE: func(cmd *cobra.Command, args []string) error {
			fs := cmd.Flags()
			fs.StringVarP(&opts.StarConfigFilePath, "config", "c", "", "config file for star")
			fs.BoolVarP(&opts.Server, "server", "s", true, "server status, true:server, false: client")
			return nil
		},

		RunE: func(cmd *cobra.Command, args []string) error {
			return runUp(&opts)
		},
	}

	return cmd
}

//runUp run a star up
func runUp(opts *upOptions) error {
	if opts.StarConfigFilePath == "" {
		opts.Server = false
	}
	tuntap, err := device.New()
	if err != nil {
		return err
	}
	fmt.Println("Create tap success", tuntap)
	return nil
}
