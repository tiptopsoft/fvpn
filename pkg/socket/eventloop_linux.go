package socket

import (
	"fmt"
	"github.com/interstellar-cloud/star/pkg/log"
	"syscall"
)

func NewEventLoop(socket *Socket) (*EventLoop, error) {
	epfd, err := syscall.EpollCreate(0)
	if err != nil {
		return nil, fmt.Errorf("create epoll fd failed: (%v)", err)
	}

	return &EventLoop{
		EpollFileDescriptor:  epfd,
		SocketFileDescriptor: socket.FileDescriptor,
	}, nil

}

func (eventLoop *EventLoop) TapFd(fd int) error {
	var event syscall.EpollEvent
	event.Events = syscall.EPOLLIN
	event.Fd = int32(fd)

	//join epfd
	if err := syscall.EpollCtl(eventLoop.EpollFileDescriptor, syscall.EPOLL_CTL_ADD, fd, &event); err != nil {
		return fmt.Errorf("add server fd to epoll failed: (%v)", err)
	}

	eventLoop.TapFileDescriptor = fd
	return nil
}

type Executor interface {
	Execute(socket Socket) error
}

func (eventLoop *EventLoop) EventLoop() {

	for {
		nevents, err := syscall.EpollWait(eventLoop.EpollFileDescriptor, eventLoop.events[:], -1)
		if err != nil {
			log.Logger.Errorf("epoll wait: (%v)", err)
		}

		for ev := 0; ev < nevents; ev++ {
			fd := eventLoop.events[ev].Fd
			var e Executor
			if int(eventLoop.events[ev].Fd) == eventLoop.SocketFileDescriptor {
				e = TapExecutor{}
			}

			if int(eventLoop.events[ev].Fd) == eventLoop.TapFileDescriptor {
				e = TapExecutor{}
			}

			if err := e.Execute(Socket{int(fd)}); err != nil {
				log.Logger.Errorf("executor fd failed. (%v)", err)
			}
		}

	}
}
