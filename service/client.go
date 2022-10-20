package service

import (
	"fmt"
	"github.com/interstellar-cloud/star/option"
	"net"
	"time"
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
	conn, err := net.DialTimeout("tcp", address, time.Second*30)
	if err != nil {
		return nil, err
	}

	return conn, nil
}
