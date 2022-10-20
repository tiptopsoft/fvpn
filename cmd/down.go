package cmd

import (
	"errors"
	"github.com/interstellar-cloud/star/device"
	"github.com/interstellar-cloud/star/option"
	"github.com/spf13/cobra"
)

type downOptions struct {
	option.StarConfig
}

func downCmd() *cobra.Command {
	var opts downOptions
	cmd := &cobra.Command{
		Use:   "down",
		Short: "remove a device",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			if len(args) == 0 {
				return errors.New("should provide at least one name of dev")
			}

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			opts.Name = args[0]
			return runDown(&opts)
		},
	}

	return cmd
}

func runDown(opts *downOptions) error {
	return device.Remove(&opts.StarConfig)
}
