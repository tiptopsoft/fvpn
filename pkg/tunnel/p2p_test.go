package tunnel

import (
	"fmt"
	"github.com/ccding/go-stun/stun"
	reuse "github.com/libp2p/go-reuseport"
	"net"
	"testing"
)

func TestInit(t *testing.T) {
	localPort := 49557
	//conn, _ := reuse.Dial("udp", fmt.Sprintf(":%d", localPort), "stun.miwifi.com:3478")
	conn, err := reuse.ListenPacket("udp", fmt.Sprintf(":%d", localPort))

	client := stun.NewClientWithConnection(conn)
	client.SetServerAddr("stun.miwifi.com:3478")
	addr := conn.LocalAddr().(*net.UDPAddr)
	fmt.Println("localhost: ", addr)
	nat, host, err := client.Discover()

	if err != nil {
		panic(err)
	}

	fmt.Println("got response:", nat, host.IP(), host.Port())

}
