package origin

import (
	"github.com/topcloudz/fvpn/pkg/socket"
	"golang.org/x/sys/unix"
	"time"
)

// import (
//
//	"fmt"
//
// )
func main() {
	//	//srcAddr := &net.UDPAddr{
	//	//	IP:   net.IPv4zero,
	//	//	Port: 6061,
	//	//}
	//	//
	//	//destAddr1 := &net.UDPAddr{
	//	//	IP:   net.ParseIP("211.125.225.186"),
	//	//	Port: 4000,
	//	//}
	//	//
	//	//destAddr2 := &net.UDPAddr{
	//	//	IP:   net.ParseIP("81.70.36.156"),
	//	//	Port: 4000,
	//	//}
	//
	//	//sock := socket.NewSocket()
	//	//destAddr1 := unix.SockaddrInet4{
	//	//	Port: 4000,
	//	//	SourceIP: [4]byte{211, 125, 225, 186},
	//	//}
	//	//
	//	//destAddr2 := unix.SockaddrInet4{
	//	//	Port: 5000,
	//	//	SourceIP: [4]byte{81, 70, 36., 156},
	//	//}
	//	//err := sock.Connect(&destAddr1)
	//	//err = sock.Connect(&destAddr2)
	//	//if err != nil {
	//	//
	//	//	panic(err)
	//	//}
	//
	//	//conn1, err := net.DialUDP("udp", srcAddr, destAddr1)
	//	//if err != nil {
	//	//	panic(err)
	//	//}
	//	//
	//	//fmt.Println(conn1)
	//	//
	//	//conn2, err := net.DialUDP("udp", srcAddr, destAddr2)
	//	//if err != nil {
	//	//	panic(err)
	//	//}
	//
	//buff := []byte{104, 101, 108, 108, 111, 44, 32, 104, 111, 108, 101, 32, 112, 117, 110, 99, 104, 105, 110, 103, 46, 46, 46, 32}
	//
	//fmt.Println(string(buff))

	sock := socket.NewSocket(6061)
	addr1 := &unix.SockaddrInet4{
		Port: 9527,
		Addr: [4]byte{211, 159, 225, 186},
	}

	addr2 := &unix.SockaddrInet4{
		Port: 9527,
		Addr: [4]byte{81, 70, 36, 156},
	}
	sock.Connect(addr1)
	sock.Connect(addr2)

	for {
		time.Sleep(time.Second * 5)
		sock.Write([]byte("hello"))
		//sock.WriteToUdp([]byte("hello"), addr2)

	}

}
