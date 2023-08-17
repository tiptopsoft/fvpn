package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/tiptopsoft/fvpn/pkg/util"
)

func logout() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "logout",
		Long: `logout out fvpn, once you logout, packet will be invalid, because fvpn check each packet`,
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			return runLogout()
		},
	}

	return cmd
}

func runLogout() error {
	localCfg, err := util.GetLocalConfig()
	if err != nil {
		return fmt.Errorf("logout failed, %v", err)
	}
	localCfg.Auth = ""
	localCfg.UserId = ""
	return util.ReplaceLocalConfig(localCfg)
}
