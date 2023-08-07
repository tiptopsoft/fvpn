package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/tiptopsoft/fvpn/pkg/util"
)

func logout() *cobra.Command {
	var opts loginOptions
	cmd := &cobra.Command{
		Use: "logout",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			return runLogout(&opts)
		},
	}

	return cmd
}

func runLogout(opts *loginOptions) error {
	localCfg, err := util.GetLocalConfig()
	if err != nil {
		return fmt.Errorf("logout failed, %v", err)
	}
	localCfg.Auth = ""
	localCfg.UserId = ""
	return util.ReplaceLocalConfig(localCfg)
}
