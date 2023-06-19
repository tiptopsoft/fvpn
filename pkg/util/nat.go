package util

import (
	"fmt"
	"github.com/ccding/go-stun/stun"
	"github.com/topcloudz/fvpn/pkg/option"
	"net"
)

var NatType uint8

func Init() {
	NatType = checkNatType()
}

// CheckNatType
func checkNatType() uint8 {

	conn, _ := net.ListenUDP("udp", nil)
	addr := conn.LocalAddr().(*net.UDPAddr)
	fmt.Println(addr.Port)
	client := stun.NewClientWithConnection(conn)
	client.SetServerAddr("stun.miwifi.com:3478")
	//client.SetServerAddr("101.43.97.112:3478")
	nat, host, err := client.Discover()

	if err != nil {
		fmt.Println(err)
		return 0
	}

	fmt.Println(nat, host)
	if nat.String() == "Symmetric NAT" {
		return option.SymmetricNAT
	}

	return option.RestrictNat
}
