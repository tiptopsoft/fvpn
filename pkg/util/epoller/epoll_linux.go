package epoller

import (
	"fmt"
	"github.com/interstellar-cloud/star/pkg/util/log"
	"github.com/interstellar-cloud/star/pkg/util/option"
	"github.com/interstellar-cloud/star/pkg/util/socket"
	"github.com/interstellar-cloud/star/pkg/util/socket/executor"
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
	sockets  map[int]*socket.Socket
}

func NewEventLoop() (*EventLoop, error) {
	epfd, err := unix.EpollCreate1(0)
	if err != nil {
		return nil, err
	}
	return &EventLoop{
		epfd:    epfd,
		sockets: make(map[int]*socket.Socket, 1),
	}, nil
}

func (eventLoop EventLoop) AddFd(skt socket.Socket) error {
	var event unix.EpollEvent
	var e error
	event.Events = unix.EPOLLIN
	//fd := socket.SocketFD(conn)
	log.Logger.Infof("Add fd: %d", skt.Fd)
	event.Fd = int32(skt.Fd)
	if eventLoop.Protocol == option.UDP {
		eventLoop.sockets[skt.Fd] = &skt
	}

	if e = unix.EpollCtl(eventLoop.epfd, unix.EPOLL_CTL_ADD, skt.Fd, &event); e != nil {
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
			var skt *socket.Socket
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
