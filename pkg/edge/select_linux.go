package edge

import (
	"github.com/interstellar-cloud/star/pkg/socket"
	"github.com/interstellar-cloud/star/pkg/tuntap"
	"github.com/interstellar-cloud/star/pkg/util/log"
	"github.com/interstellar-cloud/star/pkg/util/option"
	"golang.org/x/sys/unix"
)

type EventLoop struct {
	Tap *tuntap.Tuntap
}

func NewEventLoop(tap *tuntap.Tuntap) *EventLoop {
	return &EventLoop{Tap: tap}
}

func (eventLoop EventLoop) eventLoop(netFd, tapFd int) {
	fdMap := make(map[int]socket.Socket)
	fdMap[netFd] = socket.Socket{FileDescriptor: netFd}
	fdMap[tapFd] = socket.Socket{FileDescriptor: tapFd}
	for {
		var FdSet unix.FdSet
		var maxFd int
		if netFd > tapFd {
			maxFd = netFd
		} else {
			maxFd = tapFd
		}
		FdSet.Zero()
		FdSet.Set(netFd)
		FdSet.Set(tapFd)

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
			log.Logger.Info("tap fd start working..")
			s = fdMap[tapFd]
			executor = TapExecutor{
				Name:   eventLoop.Tap.Name,
				Socket: s,
			}
		}

		if FdSet.IsSet(netFd) {
			s = fdMap[netFd]
			executor = EdgeExecutor{
				Protocol: option.UDP,
				Tap:      eventLoop.Tap,
			}
		}

		if s.FileDescriptor != 0 {
			if err := executor.Execute(s); err != nil {
				log.Logger.Errorf("executor execute faile: (%v)", err.Error())
			}
		}

	}
}
