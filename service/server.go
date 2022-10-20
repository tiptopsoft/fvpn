package service

import (
	"errors"
	"fmt"
	"net"
)

type Reader interface {
}

var (
	Unknown = errors.New("unknown")
)

func Listen() error {
	listener, err := net.Listen("tcp", ":3000")
	if err != nil {
		return err
	}

	for {
		handle(listener)
	}

}

func handle(listener net.Listener) ([]byte, error) {

	conn, err := listener.Accept()
	defer conn.Close()
	if err != nil {
		panic(err)
	}

	buf := make([]byte, 128)
	_, err = conn.Read(buf)
	if err != nil {
		panic(err)
	}

	fmt.Println("Recived: ", buf)

	return nil, Unknown
}
