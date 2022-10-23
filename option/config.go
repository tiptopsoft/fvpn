package option

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
}

type Star struct {
	Name string
	IP   string
	Mask string
	Mode int //0 tun 1 tap
}

type Config struct {
	Listen   string
	User     string
	Name     string
	Password string
}
