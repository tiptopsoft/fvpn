package main

import (
	"fmt"
	"net"
)

func main() {

	b := []byte{192, 0, 168, 1}

	ip := net.IPv4(b[0], b[1], b[2], b[3])
	fmt.Println(ip.String())
}
