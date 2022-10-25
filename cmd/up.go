package cmd

import (
	"fmt"
	"github.com/interstellar-cloud/star/device"
	"github.com/interstellar-cloud/star/option"
	"github.com/spf13/cobra"
	"net"
	"sync"
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

	var t string
	fs.StringVarP(&t, "type", "t", "udp", "tunnel type 'tcp', 'udp'")
	if t == "tcp" {
		opts.Type = option.TCP
	} else if t == "udp" {
		opts.Type = option.UDP
	}

	return cmd
}

//runUp run a star up
func runUp(opts *upOptions) error {
	fmt.Println(fmt.Sprintf("protocol type: %d, tcp: %d, upd: %d", opts.Type, option.TCP, option.UDP))
	tun, err := device.New(&opts.StarConfig, device.TUN)
	if err != nil {
		return err
	}
	fmt.Println("Create tap success", tun)

	var netfd net.Conn
	//启动一个server
	s := &device.StarTunnel{
		Tun:  tun,
		Type: opts.Type,
	}
	if opts.Server {
		s.Serve = true
		netfd, err = s.Listen()
	}

	//是client
	if opts.StarConfig.MoonIP != "" {
		if netfd, err = s.Dial(&opts.StarConfig); err != nil {
			return err
		}
	}

	tap2net := 0
	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		defer wg.Done()
		s.Client(tap2net, netfd, tun)
	}()

	go func() {
		defer wg.Done()
		s.Server(netfd, tun)
	}()

	wg.Wait()

	return nil
}
