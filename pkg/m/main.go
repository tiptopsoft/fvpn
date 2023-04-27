package main

import (
	"fmt"
	"github.com/topcloudz/fvpn/pkg/tuntap"
)

func main() {

	//tuntap.GetTuntap("c04d6b84fd4fc978")

	ifr := tuntap.Ifreq{}

	net := "fvpnc04d6b84fd4fc978"

	copy(ifr.Name[:], "fvpnc04d6b84fd4fc978")
	fmt.Println(string(ifr.Name[:]))
	fmt.Println(net[0:16])
}
