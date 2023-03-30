package cmd

import "github.com/spf13/cobra"

type logoutOptions struct{}

func logoutCmd() *cobra.Command {
	var opts logoutOptions
	var cmd = &cobra.Command{
		Use:          "logout",
		SilenceUsage: true,
		Short:        "user logout",
		Long:         `user logout`,

		PreRunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},

		RunE: func(cmd *cobra.Command, args []string) error {
			return runLogout(&opts)
		},
	}

	return cmd
}

func runLogout(opts *logoutOptions) error {
	//TODO to be implemented
	return nil
}
