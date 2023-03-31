package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:          "fvpn [command]",
	SilenceUsage: true,
	Short:        "Start a fvpn, which can visit you own network.",
	Long:         `Start a fvpn, use which can visit private net.`,
}

func Execute() {
	rootCmd.AddCommand(EdgeCmd(), RegCmd())
	if err := rootCmd.Execute(); err != nil {
		os.Exit(-1)
	}
}
