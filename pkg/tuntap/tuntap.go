package tuntap

import (
	"net"
	"os"
)

// Tuntap a tuntap for net
type Tuntap struct {
	file    *os.File
	Fd      int
	Name    string
	Mode    Mode
	MacAddr net.HardwareAddr
}

type Mode int

const (
	TUN Mode = iota
	TAP
)

var NamePrefix string = "fvpn"

type Ifreq struct {
	Name  [16]byte
	Flags uint16
}

func (t *Tuntap) Write(buff []byte) (int, error) {
	return t.file.Write(buff)
}

func (t *Tuntap) Read(buff []byte) (int, error) {
	return t.file.Read(buff)
}
