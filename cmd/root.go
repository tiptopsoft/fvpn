package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:          "fvpn [command]",
	SilenceUsage: true,
	Short:        "fvpn can let you join/leave a private network, compose network",
	Long:         `fvpn can let you join/leave a private network, compose our network, build node mesh and so on`,
}

func Execute() {
	rootCmd.AddCommand(nodeCmd(), registryCmd(), joinCmd(), loginCmd(), logout(), leaveCmd(), statusCmd(), stopCmd())
	if err := rootCmd.Execute(); err != nil {
		os.Exit(-1)
	}
}
