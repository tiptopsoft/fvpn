package util

import (
	"fmt"
	"github.com/ccding/go-stun/stun"
	"github.com/topcloudz/fvpn/pkg/option"
)

var NatType uint8

func init() {
	NatType = checkNatType()
}

// CheckNatType
func checkNatType() uint8 {
	client := stun.NewClient()
	client.SetServerAddr("stun.miwifi.com:3478")
	nat, host, err := client.Discover()
	if err != nil {
		return 0
	}

	fmt.Println(nat, host)
	if nat.String() == "Symmetric NAT" {
		return option.SymmetricNAT
	}

	return option.RestrictNat
}
