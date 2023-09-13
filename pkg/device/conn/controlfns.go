package conn

import (
	"net"
	"syscall"
)

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
