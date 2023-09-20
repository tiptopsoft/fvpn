package conn

import (
	"net"
	"syscall"
)

// UDP socket read/write buffer size (7MB). The value of 7MB is chosen as it is
// the max supported by a default configuration of macOS. Some platforms will
// silently clamp the value to other maximums, such as linux clamping to
// net.core.{r,w}mem_max (see _linux.go for additional implementation that works
// around this limitation)
const socketBufferSize = 7 << 20

type fn func(network, address string, c syscall.RawConn) error

var fns []fn

func ListenConfig() *net.ListenConfig {
	return &net.ListenConfig{
		Control: func(network, address string, c syscall.RawConn) error {
			for _, f := range fns {
				err := f(network, address, c)
				if err != nil {
					return err
				}
			}

			return nil
		},
	}
}
