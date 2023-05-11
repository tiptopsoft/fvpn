package main

import (
	"fmt"
	"github.com/spf13/cobra"
)

var (
	// Used for flags.
	cfgFile     string
	userLicense string
)

func rootCmd() *cobra.Command {
	var author string
	rootCmd := &cobra.Command{
		Use:   "cli",
		Short: "A generator for Cobra based Applications",
		Long: `Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
		Run: func(cmd *cobra.Command, args []string) {
			run(author)
		},
	}

	rootCmd.PersistentFlags().StringVar(&author, "author", "YOUR NAME", "Author name for copyright attribution")
	return rootCmd
}

// Execute executes the root command.
func Execute() error {
	cmd := rootCmd()
	return cmd.Execute()
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number of Hugo",
	Long:  `All software has versions. This is Hugo's`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Hugo Static Site Generator v0.9 -- HEAD")
	},
}

func main() {
	Execute()
}

func run(name string) error {
	fmt.Println("author:", name)
	return nil
}
