package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:          "fvpn [command]",
	SilenceUsage: true,
	Short:        "use fvpn, can start a node or a registry, join or leave a network",
	Long:         `use fvpn tools, can let you logout our service, join network and so on`,
}

func Execute() {
	rootCmd.AddCommand(EdgeCmd(), RegCmd(), joinCmd(), loginCmd(), logout(), leaveCmd())
	if err := rootCmd.Execute(); err != nil {
		os.Exit(-1)
	}
}
