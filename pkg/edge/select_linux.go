package edge

import (
	"github.com/interstellar-cloud/star/pkg/socket"
	"github.com/interstellar-cloud/star/pkg/tuntap"
	"github.com/interstellar-cloud/star/pkg/util/log"
	"github.com/interstellar-cloud/star/pkg/util/option"
	"golang.org/x/sys/unix"
	"sync"
)

var fdMap sync.Map
var FdSet unix.FdSet

type EventLoop struct {
	Tap *tuntap.Tuntap
}

func NewEventLoop(tap *tuntap.Tuntap) *EventLoop {
	return &EventLoop{Tap: tap}
}

func (eventLoop EventLoop) eventLoop(netFd, tapFd int) {

	fdMap.Store(netFd, socket.Socket{FileDescriptor: netFd})
	fdMap.Store(tapFd, socket.Socket{FileDescriptor: tapFd})
	for {
		var maxFd int
		if netFd > tapFd {
			maxFd = netFd
		} else {
			maxFd = tapFd
		}
		FdSet.Zero()
		FdSet.Set(netFd)
		FdSet.Set(tapFd)

		for {
			ret, err := unix.Select(maxFd+1, &FdSet, nil, nil, nil)
			if ret < 0 && err == unix.EINTR {
				continue
			}
			var s socket.Socket
			var executor socket.Executor
			if err != nil {
				panic(err)
			}

			if FdSet.IsSet(tapFd) {
				sAny, _ := fdMap.Load(tapFd)
				s = sAny.(socket.Socket)
				executor = TapExecutor{
					Name:   eventLoop.Tap.Name,
					Socket: socket.Socket{FileDescriptor: netFd},
				}
			}

			if FdSet.IsSet(netFd) {
				sAny, _ := fdMap.Load(netFd)
				s = sAny.(socket.Socket)
				executor = EdgeExecutor{
					Protocol: option.UDP,
					Tap:      eventLoop.Tap,
				}
			}

			if s.FileDescriptor != 0 {
				if err := executor.Execute(s); err != nil {
					log.Logger.Errorf("executor execute faile: (%v)", err)
				}
			}

		}
	}
}

func AddFd(fd int) {
	FdSet.Set(fd)
}
