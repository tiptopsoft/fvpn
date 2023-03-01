package epoller

import (
	"fmt"
	"github.com/interstellar-cloud/star/pkg/executor"
	"github.com/interstellar-cloud/star/pkg/log"
	"github.com/interstellar-cloud/star/pkg/option"
	socket "github.com/interstellar-cloud/star/pkg/socket"
	"golang.org/x/sys/unix"
	"os"
	"syscall"
)

const (
	MaxEpollEvents = 32
)

type EventLoop struct {
	epfd           int
	events         [MaxEpollEvents]syscall.EpollEvent
	fileDescriptor int
	//*registry.RegStar
	Protocol option.Protocol
	sockets  map[int]*socket.Interface
}

func NewEventLoop() (*EventLoop, error) {
	epfd, err := unix.EpollCreate1(0)
	if err != nil {
		return nil, err
	}
	return &EventLoop{
		epfd:    epfd,
		sockets: make(map[int]*socket.Interface, 1),
	}, nil
}

func (eventLoop EventLoop) AddFd(skt socket.Interface) error {
	var event unix.EpollEvent
	var e error
	event.Events = unix.EPOLLIN
	//fd := socket.SocketFD(conn)
	sktNew := skt.(socket.Socket)
	event.Fd = int32(sktNew.Fd)
	if eventLoop.Protocol == option.UDP {
		eventLoop.sockets[sktNew.Fd] = &skt
	}

	if e = unix.EpollCtl(eventLoop.epfd, unix.EPOLL_CTL_ADD, sktNew.Fd, &event); e != nil {
		fmt.Println("epoll_ctl: ", e)
		os.Exit(-1)
	}

	return nil
}

func (eventLoop *EventLoop) EventLoop(executor executor.Executor) {
	for {
		events := eventLoop.events
		nevents, e := syscall.EpollWait(eventLoop.epfd, eventLoop.events[:], -1)
		if e != nil {
			fmt.Println("epoll_wait: ", e)
			continue
		}

		for i := 0; i < nevents; i++ {
			fd := events[i].Fd
			var skt *socket.Interface
			skt = eventLoop.sockets[int(fd)]

			if skt != nil {
				go func() {
					if err := executor.Execute(*skt); err != nil {
						log.Logger.Errorf("executor exec failed. (%v)", err)
					}
				}()
			}
		}

	}
}
