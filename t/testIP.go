package main

import (
	"fmt"
	"net/netip"
)

func main() {
	addr, err := netip.ParseAddr("192.168.0.1")
	fmt.Println(addr, err)
}
