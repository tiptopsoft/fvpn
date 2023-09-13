package conn

import (
	"fmt"
	"golang.org/x/sys/unix"
	"runtime"
	"syscall"
)

// set opt for linux
func init() {

	fns = append(fns, func(network, address string, c syscall.RawConn) error {
		var err error
		switch network {
		case "udp4":
			if runtime.GOOS != "android" {
				c.Control(func(fd uintptr) {
					err = unix.SetsockoptInt(int(fd), unix.IPPROTO_IP, unix.IP_PKTINFO, 1)
				})
			}
		case "udp6":
			c.Control(func(fd uintptr) {
				if runtime.GOOS != "android" {
					err = unix.SetsockoptInt(int(fd), unix.IPPROTO_IPV6, unix.IPV6_RECVPKTINFO, 1)
					if err != nil {
						return
					}
				}
				err = unix.SetsockoptInt(int(fd), unix.IPPROTO_IPV6, unix.IPV6_V6ONLY, 1)
			})
		default:
			err = fmt.Errorf("unhandled network: %s: %w", network, unix.EINVAL)
		}
		return err
	})
}
