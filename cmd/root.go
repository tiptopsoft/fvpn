package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:          "fvpn [command]",
	SilenceUsage: true,
	Short:        "fvpn can let you join/leave a private network, compose network",
	Long:         `fvpn can let you join/leave a private network, compose network`,
}

func Execute() {
	rootCmd.AddCommand(EdgeCmd(), RegCmd(), joinCmd(), loginCmd(), logout(), leaveCmd())
	if err := rootCmd.Execute(); err != nil {
		os.Exit(-1)
	}
}
