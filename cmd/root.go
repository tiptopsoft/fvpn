package cmd

import (
	"github.com/spf13/cobra"
	"os"
)

var rootCmd = &cobra.Command{
	Use:          "star [command]",
	SilenceUsage: true,
	Short:        "Start a star, use which can visit private net.",
	Long:         `Start a star, use which can visit private net.`,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(-1)
	}
}
