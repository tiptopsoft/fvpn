package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:          "star [command]",
	SilenceUsage: true,
	Short:        "Start a edge, which can visit you own network.",
	Long:         `Start a edge, use which can visit private net.`,
}

func Execute() {
	rootCmd.AddCommand(edgeCmd(), regCmd())
	if err := rootCmd.Execute(); err != nil {
		os.Exit(-1)
	}
}
