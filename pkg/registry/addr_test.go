package registry

import (
	"fmt"
	"github.com/interstellar-cloud/star/pkg/packet"
	ack2 "github.com/interstellar-cloud/star/pkg/packet/register/ack"
	"net"
	"testing"
)

func TestGenerateIP(t *testing.T) {
	//ip1 := string2Long("192.168.0.1")
	//fmt.Println(ip1)
	//fmt.Println(strconv.FormatInt(36, 10))
	//fmt.Println(GenerateIP(ip1))

	ip, _, _ := net.ParseCIDR("255.255.255.0/24")

	fmt.Println(ip.String())
	ack := ack2.RegPacketAck{}
	mask := make([]byte, 16)
	packet.DecodeBytes(&mask, ip, 0)
	ack.Mask = mask

	fmt.Println(ack.Mask.String())
}
