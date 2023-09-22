package conn

import (
	"golang.org/x/sys/windows"
	"syscall"
)

func init() {
	fns = append(fns, // Attempt to set the socket buffer size beyond net.core.{r,w}mem_max by
		// using SO_*BUFFORCE. This requires CAP_NET_ADMIN, and is allowed here to
		// fail silently - the result of failure is lower performance on very fast
		// links or high latency links.
		func(network, address string, c syscall.RawConn) error {
			return c.Control(func(fd uintptr) {
				_ = windows.SetsockoptInt(windows.Handle(fd), windows.SOL_SOCKET, windows.SO_RCVBUF, socketBufferSize)
				_ = windows.SetsockoptInt(windows.Handle(fd), windows.SOL_SOCKET, windows.SO_SNDBUF, socketBufferSize)
			})
		})
}
