package main

import (
	"fmt"
	"github.com/interstellar-cloud/star/pkg/util/tuntap"
	"github.com/spf13/cobra"
	"golang.org/x/sys/unix"
	"os"
)

func main() {
	Execute()
}

func runMain() {
	tap, err := tuntap.New(tuntap.TAP)
	if err != nil {
		panic(err)
	}

	for {
		var FdSet unix.FdSet
		var maxFd int
		tapFd := int(tap.Fd)
		maxFd = int(tap.Fd)
		FdSet.Zero()
		FdSet.Set(tapFd)

		ret, err := unix.Select(maxFd+1, &FdSet, nil, nil, nil)
		if ret < 0 && err == unix.EINTR {
			continue
		}

		if err != nil {
			panic(err)
		}

		if FdSet.IsSet(tapFd) {
			b := make([]byte, 1024)
			n, err := tap.Socket.Read(b)
			if err != nil {
				fmt.Println(err)
			}

			fmt.Println(fmt.Sprintf("Read %d bytes from device %s", n, tap.Name))
		}
	}

}

var rootCmd = &cobra.Command{
	Use:          "main",
	SilenceUsage: true,
	Short:        "Start a edge, use which can visit private net.",
	Long:         `Start a edge, use which can visit private net.`,
	Run: func(cmd *cobra.Command, args []string) {
		runMain()
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(-1)
	}
}
