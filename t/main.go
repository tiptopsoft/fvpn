package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"net"
	"strconv"
)

//ip到数字
func ip2Long(ip string) uint32 {
	var long uint32
	binary.Read(bytes.NewBuffer(net.ParseIP(ip).To4()), binary.BigEndian, &long)
	return long
}

//数字到IP
func backtoIP4(ipInt int64) string {
	// need to do two bit shifting and “0xff” masking
	b0 := strconv.FormatInt((ipInt>>24)&0xff, 10)
	b1 := strconv.FormatInt((ipInt>>16)&0xff, 10)
	b2 := strconv.FormatInt((ipInt>>8)&0xff, 10)
	b3 := strconv.FormatInt((ipInt & 0xff), 10)
	return b0 + "." + b1 + "." + b2 + "." + b3
}
func main() {
	result := ip2Long("98.138.253.109")
	fmt.Println(result)
	// or if you prefer the super fast way
	faster := binary.BigEndian.Uint32(net.ParseIP("98.138.253.109")[12:16])
	fmt.Println(faster)
	faster64 := int64(faster)
	fmt.Println(backtoIP4(faster64))
	ip1 := ip2Long("221.177.0.0")
	ip2 := ip2Long("221.177.7.255")
	//ip1 := ip2Long("192.168.0.0")
	//ip2 := ip2Long("192.168.0.255")
	x := ip2 - ip1
	fmt.Println(ip1, ip2, x)
	for i := ip1; i <= ip2; i++ {
		i := int64(i)
		fmt.Println(backtoIP4(i))
	}
}
