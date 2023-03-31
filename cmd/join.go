package cmd

import "github.com/spf13/cobra"

type joinOptions struct {
}

func joinCmd() *cobra.Command {
	var opts joinOptions
	var cmd = &cobra.Command{
		Use:          "join",
		SilenceUsage: true,
		Short:        "user login user username and password which registered in our website",
		Long:         `user login user username and password which registered in our website`,

		PreRunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},

		RunE: func(cmd *cobra.Command, args []string) error {
			return runJoin(&opts)
		},
	}
	//fs := cmd.Flags()
	//fs.StringVarP(&opts.StarConfigFilePath, "config", "", "", "config file for client")

	return cmd
}

func runJoin(opts *joinOptions) error {
	//TODO to be implemented
	return nil
}
