package cmd

import "github.com/spf13/cobra"

type loginOptions struct {
}

func loginCmd() *cobra.Command {
	var opts loginOptions
	var cmd = &cobra.Command{
		Use:          "login",
		SilenceUsage: true,
		Short:        "user login user username and password which registered in our website",
		Long:         `user login user username and password which registered in our website`,

		PreRunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},

		RunE: func(cmd *cobra.Command, args []string) error {
			return runLogin(&opts)
		},
	}
	//fs := cmd.Flags()
	//fs.StringVarP(&opts.StarConfigFilePath, "config", "", "", "config file for edge")

	return cmd
}

func runLogin(opts *loginOptions) error {
	//TODO to be implemented
	return nil
}
