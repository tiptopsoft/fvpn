package main

import (
	"fmt"
	"github.com/topcloudz/fvpn/pkg/packet/register"
	"net"
)

func main() {
	mac, err := net.ParseMAC("00:00:00:00:fe:80:00:00:00:00:00:00:02:00:5e:10:00:00:00:01")
	if err != nil {
		panic(err)
	}

	ip := net.IP{192, 168, 0, 1}

	reg := register.NewPacket("c04d6b84fd4fc978", mac, ip)
	bs, err := register.Encode(reg)
	if err != nil {
		panic(err)
	}

	fmt.Println(bs)
	fmt.Println(len(bs))

	r, _ := register.Decode(bs)
	fmt.Println(r)
}
