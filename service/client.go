package service

import (
	"fmt"
	"github.com/interstellar-cloud/star/option"
	"net"
)

const DefaultPort = 3000

type Client struct {
}

func Conn(opts *option.StarConfig) (net.Conn, error) {
	if opts.Port == 0 {
		opts.Port = DefaultPort
	}
	address := fmt.Sprintf("%s:%d", opts.MoonIP, opts.Port)
	fmt.Println("connect to:", address)
	conn, err := net.Dial("tcp", address)
	if err != nil {
		return nil, err
	}

	return conn, nil
}
