package cmd

import (
	"fmt"
	"github.com/interstellar-cloud/star/device"
	"github.com/interstellar-cloud/star/option"
	"github.com/interstellar-cloud/star/service"
	"github.com/spf13/cobra"
	"io"
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

	var netfd net.Conn
	//启动一个server
	if opts.Server {
		s := &service.Server{
			Tun: tun,
		}
		netfd, err = s.Listen()
		if err != nil {
			return err
		}

	}

	//是client
	if opts.StarConfig.MoonIP != "" {
		netfd, err = service.Conn(&opts.StarConfig)
		if err != nil {
			return err
		}
	}

	tap2net := 0
	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		defer wg.Done()
		client(tap2net, netfd, tun)
	}()

	go func() {
		defer wg.Done()
		server(netfd, tun)
	}()

	wg.Wait()

	return nil
}

func client(tap2net int, netfd io.ReadWriteCloser, tun *device.Tuntap) {

	for {
		var buf [2000]byte
		n, err := tun.Read(buf[:])
		if err == io.EOF {
			continue
		}
		if err != nil {
			panic(err)
		}

		tap2net++
		fmt.Println(fmt.Printf("tap2net:%d, tun received %d byte from %s: ", tap2net, n, tun.Name))

		/* write packet */
		n, err = netfd.Write(buf[:n])
		if err != nil {
			panic(err)
		}

		fmt.Println(fmt.Printf("tap2net:%d,write %d byte to network", tap2net, n))
	}

}

func server(netfd io.ReadWriteCloser, tun *device.Tuntap) {
	for {
		buf := make([]byte, 2000)
		n, err := netfd.Read(buf)
		if err == io.EOF {
			continue
		}
		fmt.Println(fmt.Printf("Recevied %d byte from net", n))
		//write to tap
		_, err = tun.Write(buf[:n])
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println(fmt.Printf("write %d byte to tap %s", n, tun.Name))
	}

}
