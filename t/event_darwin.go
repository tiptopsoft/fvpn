package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	"strconv"
	"strings"
	"syscall"
)

type Socket struct {
	FileDescriptor int
}

func (socket Socket) Read(bytes []byte) (int, error) {
	if len(bytes) == 0 {
		return 0, nil
	}

	numBytesRead, err := syscall.Read(socket.FileDescriptor, bytes)
	if err != nil {
		numBytesRead = 0
	}

	return numBytesRead, err
}

func (socket Socket) Write(bytes []byte) (int, error) {
	numBytesWritten, err := syscall.Write(socket.FileDescriptor, bytes)
	if err != nil {
		numBytesWritten = 0
	}

	return numBytesWritten, err
}

func (socket Socket) Close() error {
	return syscall.Close(socket.FileDescriptor)
}

func (socket *Socket) String() string {
	return strconv.Itoa(socket.FileDescriptor)
}

func Listen(ip string, port int) (*Socket, error) {
	socket := &Socket{}

	socketFileDescriptor, err := syscall.Socket(syscall.AF_INET, syscall.SOCK_STREAM, 0)
	if err != nil {
		return nil, fmt.Errorf("failed to create socket (%v)", err)
	}

	socket.FileDescriptor = socketFileDescriptor
	socketAddress := &syscall.SockaddrInet4{Port: port}
	copy(socketAddress.Addr[:], net.ParseIP(ip))

	if err = syscall.Bind(socket.FileDescriptor, socketAddress); err != nil {
		return nil, fmt.Errorf("failed to bind socket (%v)", err)
	}

	if err = syscall.Listen(socket.FileDescriptor, syscall.SOMAXCONN); err != nil {
		return nil, fmt.Errorf("failed to listen on socket (%v)", err)
	}

	return socket, nil

}

type EventLoop struct {
	KqueueFileDescriptor int
	SocketFileDescriptor int
}

func NewEventLoop(s *Socket) (*EventLoop, error) {
	kqueue, err := syscall.Kqueue()
	if err != nil {
		return nil, fmt.Errorf("failed to create kqueue file descriptor (%v)", err)
	}

	changeEvent := syscall.Kevent_t{
		Ident:  uint64(s.FileDescriptor),
		Filter: syscall.EVFILT_READ,
		Flags:  syscall.EV_ADD | syscall.EV_ENABLE,
		Fflags: 0,
		Data:   0,
		Udata:  nil,
	}

	changeEventRegistered, err := syscall.Kevent(kqueue, []syscall.Kevent_t{changeEvent}, nil, nil)

	if err != nil || changeEventRegistered == -1 {
		return nil, fmt.Errorf("failed to register change eventloop (%v)", err)
	}

	return &EventLoop{
		KqueueFileDescriptor: kqueue,
		SocketFileDescriptor: s.FileDescriptor,
	}, nil

}

type Handler func(socket *Socket)

func (eventLoop EventLoop) Handle(handler Handler) {

	for {
		newEvents := make([]syscall.Kevent_t, 10)
		numNewEvents, err := syscall.Kevent(
			eventLoop.KqueueFileDescriptor,
			nil,
			newEvents,
			nil,
		)
		if err != nil {
			continue
		}

		for i := 0; i < numNewEvents; i++ {
			currentEvent := newEvents[i]
			eventFileDescriptor := int(currentEvent.Ident)

			if currentEvent.Flags&syscall.EV_EOF != 0 {
				// client closing connection
				syscall.Close(eventFileDescriptor)
			} else if eventFileDescriptor == eventLoop.SocketFileDescriptor {
				// new incoming connection
				socketConnection, _, err := syscall.Accept(eventFileDescriptor)
				if err != nil {
					continue
				}

				socketEvent := syscall.Kevent_t{
					Ident:  uint64(socketConnection),
					Filter: syscall.EVFILT_READ,
					Flags:  syscall.EV_ADD,
					Fflags: 0,
					Data:   0,
					Udata:  nil,
				}
				socketEventRegistered, err := syscall.Kevent(
					eventLoop.KqueueFileDescriptor,
					[]syscall.Kevent_t{socketEvent},
					nil,
					nil,
				)
				if err != nil || socketEventRegistered == -1 {
					continue
				}
			} else if currentEvent.Filter&syscall.EVFILT_READ != 0 {
				// data available -> forward to handler
				handler(&Socket{
					FileDescriptor: int(eventFileDescriptor),
				})
			}

			// ignore all other events
		}
	}
}

func main() {

	s, err := Listen("127.0.0.1", 8080)
	if err != nil {
		log.Println("Failed to create socket:", err)
		os.Exit(1)
	}

	eventLoop, err := NewEventLoop(s)
	if err != nil {
		log.Println("Failed to create eventloop loop:", err)
		os.Exit(1)
	}

	log.Println("Server started. Waiting for incoming connections. ^C to exit.")

	eventLoop.Handle(func(s *Socket) {
		reader := bufio.NewReader(s)
		for {
			line, err := reader.ReadString('\n')
			if err != nil || strings.TrimSpace(line) == "" {
				break
			}
			s.Write([]byte(line))
		}
		s.Close()
	})
}
