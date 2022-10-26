package option

import (
	"context"
	"fmt"
	"math/rand"
	"net"
)

const (
	TCP = iota
	UDP
)

// StarConfig conf for running a star up
type StarConfig struct {
	Star
	MoonIP string // service for moon server
	Port   int    // default port is 3000
	Server bool   // server or client, true: server
	Mac    string // like "07:00:10:24:55:42"
}

type SuperStar struct {
	Listen int
}

type Star struct {
	Name string
	IP   string
	Mask string
	Mode int //0 tun 1 tap
}

type Config struct {
	Listen string `mapstructure:"listen"`
	Mysql  Mysql  `mapstructure:"mysql"`
}

type Mysql struct {
	User     string `mapstructure:"user"`
	Url      string `mapstructure:"url"`
	Password string `mapstructure:"password"`
	Name     string `mapstructure:"name"`
}

// PeerInfo star info in a machine
// containes a socket bind it, can register to superNode or to star
type PeerInfo struct {
	Mac     string
	IP      string
	Port    int
	Fd      net.Conn
	ExtIP   string
	ExtPort int
}

// RandMac rand gen a mac
func RandMac(ctx context.Context) (string, error) {
	buf := make([]byte, 6)
	_, err := rand.Read(buf)
	if err != nil {
		return "", err
	}
	buf[0] |= 2
	mac := fmt.Sprintf("%02x:%02x:%02x:%02x:%02x:%02x", buf[0], buf[1], buf[2], buf[3], buf[4], buf[5])

	return mac, nil
}
