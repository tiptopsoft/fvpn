package epoll

import (
	"fmt"
	"github.com/interstellar-cloud/star/pkg/log"
	"github.com/interstellar-cloud/star/pkg/option"
	"github.com/interstellar-cloud/star/pkg/socket"
	"syscall"
)

const (
	MaxEpollEvents = 32
)

type EventLoop struct {
	EpollFileDescriptor  int
	events               [MaxEpollEvents]syscall.EpollEvent
	SocketFileDescriptor int //Server fd
	TapFileDescriptor    int
}

func NewEventLoop(socket *socket.Socket) (*EventLoop, error) {
	epfd, err := syscall.EpollCreate1(0)
	if err != nil {
		return nil, fmt.Errorf("create epoll fd failed: (%v)", err)
	}

	eventLoop := &EventLoop{
		EpollFileDescriptor:  epfd,
		SocketFileDescriptor: socket.FileDescriptor,
	}

	var event syscall.EpollEvent
	event.Events = syscall.EPOLLIN
	event.Fd = int32(socket.FileDescriptor)

	//join epfd
	if err := syscall.EpollCtl(eventLoop.EpollFileDescriptor, syscall.EPOLL_CTL_ADD, socket.FileDescriptor, &event); err != nil {
		return nil, fmt.Errorf("add server fd to epoll failed: (%v)", err)
	}

	return eventLoop, nil
}

func (eventLoop *EventLoop) TapFd(socket socket.Socket) error {
	var event syscall.EpollEvent
	event.Events = syscall.EPOLLIN
	event.Fd = int32(socket.FileDescriptor)

	//join epfd
	if err := syscall.EpollCtl(eventLoop.EpollFileDescriptor, syscall.EPOLL_CTL_ADD, socket.FileDescriptor, &event); err != nil {
		return fmt.Errorf("add tap fd to epoll failed: (%v)", err)
	}

	eventLoop.TapFileDescriptor = socket.FileDescriptor
	return nil
}

type Executor interface {
	Execute(socket socket.Socket) error
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
			if int(fd) == eventLoop.SocketFileDescriptor {
				e = EdgeExecutor{
					Protocol:  option.UDP,
					EventLoop: eventLoop,
				}
			}

			if int(fd) == eventLoop.TapFileDescriptor {
				e = TapExecutor{}
			}

			if err := e.Execute(socket.Socket{FileDescriptor: int(fd)}); err != nil {
				log.Logger.Errorf("executor fd failed. (%v)", err)
			}
		}

	}
}
