package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:          "fvpn [command]",
	SilenceUsage: true,
	Short:        "Start a client, which can visit you own network.",
	Long:         `Start a client, use which can visit private net.`,
}

func Execute() {
	rootCmd.AddCommand(nodeCmd(), regCmd())
	if err := rootCmd.Execute(); err != nil {
		os.Exit(-1)
	}
}
