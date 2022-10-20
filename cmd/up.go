package cmd

import (
	"fmt"
	"github.com/interstellar-cloud/star/device"
	"github.com/interstellar-cloud/star/option"
	"github.com/interstellar-cloud/star/service"
	"github.com/spf13/cobra"
)

type upOptions struct {
	option.StarConfig
	option.StarAuth

	StarConfigFilePath string
}

func upCmd() *cobra.Command {
	var opts upOptions
	var cmd = &cobra.Command{
		Use:          "up",
		SilenceUsage: true,
		Short:        "start up a star, for net proxy",
		Long:         `Start up a star, for private net proxy`,

		PreRunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},

		RunE: func(cmd *cobra.Command, args []string) error {
			return runUp(&opts)
		},
	}
	fs := cmd.Flags()
	fs.StringVarP(&opts.StarConfigFilePath, "config", "", "", "config file for star")
	fs.BoolVarP(&opts.Server, "server", "s", false, "server status, true:server, false: client")
	fs.StringVarP(&opts.IP, "ip", "", "", "star config, ip")
	fs.StringVarP(&opts.Name, "name", "i", "", "star config, tuntap name")
	fs.StringVarP(&opts.Mask, "mask", "", "", "tuntap mask")
	fs.StringVarP(&opts.MoonIP, "host", "c", "", "tun server")
	fs.IntVarP(&opts.Port, "port", "p", 3000, "tun server port")

	return cmd
}

//runUp run a star up
func runUp(opts *upOptions) error {
	fmt.Println("server", opts.Server)
	tun, err := device.New(&opts.StarConfig)
	if err != nil {
		return err
	}
	fmt.Println("Create tap success", tun)

	//启动一个server
	if opts.Server {
		return service.Listen()
	}

	//是client
	if opts.StarConfig.MoonIP != "" {
		conn, err := service.Conn(&opts.StarConfig)
		defer conn.Close()
		if err != nil {
			return err
		}

		for {
			var buf []byte
			n, err := tun.Read(buf)
			if err != nil {
				return err
			}

			fmt.Println("tun received byte: ", n, buf)

			n, err = conn.Write(buf)
			if err != nil {
				return err
			}

			fmt.Println("tun write byte:", n, buf)
		}
	}

	//Read data from tuntap
	return nil
}
