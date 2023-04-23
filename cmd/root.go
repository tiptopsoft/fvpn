package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:          "fvpn [command]",
	SilenceUsage: true,
	Short:        "Start a fvpn, will read from your network",
	Long:         `Start a fvpn, use which can visit private net.`,
}

func Execute() {
	rootCmd.AddCommand(EdgeCmd(), RegCmd(), joinCmd())
	if err := rootCmd.Execute(); err != nil {
		os.Exit(-1)
	}
}
