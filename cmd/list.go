package cmd

import "github.com/spf13/cobra"

type listOptions struct {
}

func listCmd() *cobra.Command {
	opts := listOptions{}
	cmd := &cobra.Command{
		Use:          "list",
		SilenceUsage: true,
		Short:        "list",
		Long:         `when you've login in, list will show your networkIds in pretty json'`,

		PreRunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},

		RunE: func(cmd *cobra.Command, args []string) error {
			return runList(opts)
		},
	}

	return cmd
}

func runList(options listOptions) error {

	return nil
}
