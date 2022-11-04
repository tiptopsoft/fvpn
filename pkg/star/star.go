package star

import (
	"errors"
	"fmt"
	"github.com/interstellar-cloud/star/pkg/device"
	"github.com/interstellar-cloud/star/pkg/option"
	"github.com/interstellar-cloud/star/pkg/pack"
	"io"
	"net"
)

type EdgeStar struct {
	Tap   *device.Tuntap
	Addr  *net.UDPAddr
	Type  int
	Serve bool
	Conn  net.Conn
}

func (es *EdgeStar) Start(port int) error {
	if port == 0 {
		port = int(pack.DefaultPort)
	}
	conn, err := es.listen(fmt.Sprintf(":%d", port))
	if err != nil {
		return nil
	}
	es.Conn = conn

	if err := es.register(); err != nil {
		return errors.New("register to register star failed")
	}

	//get group edgestarlist
	return nil
}

func (es *EdgeStar) listen(address string) (net.Conn, error) {
	var conn net.Conn
	var err error

	switch es.Type {
	case option.TCP:
		listener, err := net.Listen("tcp", address)
		if err != nil {
			return nil, err
		}
		conn, err = listener.Accept()
	case option.UDP:
		conn, err = net.ListenUDP("udp", &net.UDPAddr{IP: net.IPv4(0, 0, 0, 0),
			Port: int(pack.DefaultPort)})
	}

	//defer conn.Close()
	if err != nil {
		panic(err)
	}

	return conn, nil
}

// register register a edgestar to center.
func (es *EdgeStar) register() error {
	p := pack.NewPacket()
	p.Flags = pack.TAP_REGISTER
	p.TTL = pack.DefaultTTL

	mac, err := option.GetLocalMac(es.Tap.Name)
	if err != nil {
		return option.ErrGetMac
	}
	copy(p.SourceMac[:], mac[:])

	data, err := pack.Encode(p)
	if err != nil {
		return errors.New("encode packet failed")
	}

	switch es.Type {
	case option.UDP:
		if _, err := es.Conn.(*net.UDPConn).Write(data); err != nil {
			return err
		}
		break
	}
	return nil
}

// listEdgeStar get all group star from super and connect to them.
func (es *EdgeStar) listEdgeStar() ([]EdgeStar, error) {

	p := pack.NewPacket()
	p.Flags = pack.TAP_LIST_EDGE_STAR

	var err error
	data, err := pack.Encode(p)
	if err != nil {
		return nil, err
	}

	if _, err := es.Conn.Write(data); err != nil {
		return nil, err
	}

	var bs = make([]byte, pack.FRAME_SIZE)
	if _, _, err = es.Conn.(*net.UDPConn).ReadFromUDP(bs[:]); err != nil {
		return nil, err
	}

	return nil, nil
}

func (es *EdgeStar) Dial(opts *option.StarConfig) (net.Conn, error) {
	if opts.Port == 0 {
		opts.Port = int(pack.DefaultPort)
	}
	address := fmt.Sprintf("%s:%d", opts.MoonIP, opts.Port)
	fmt.Println("connect to:", address)
	var conn net.Conn
	var err error
	switch es.Type {
	case option.TCP:
		conn, err = net.Dial("tcp", address)
	case option.UDP:
		ip := net.ParseIP(opts.MoonIP)
		conn, err = net.DialUDP("udp", nil, &net.UDPAddr{IP: ip,
			Port: int(pack.DefaultPort)})
	}

	if err != nil {
		return nil, err
	}

	return conn, nil
}

func (es *EdgeStar) Client(tap2net int, netfd io.ReadWriteCloser, tun *device.Tuntap) {

	for {
		var buf [2000]byte
		n, err := tun.Read(buf[:])
		if err == io.EOF {
			continue
		}
		if err != nil {
			panic(err)
		}

		tap2net++
		fmt.Println(fmt.Printf("tap2net:%d, tun received %d byte from %s: ", tap2net, n, tun.Name))

		/* write pack */
		if es.Type == option.UDP && es.Serve {
			fmt.Println("using write 2 udp")
			netfd.(*net.UDPConn).WriteToUDP(buf[:], es.Addr)
		} else {
			fmt.Println("using write udp")
			n, err = netfd.Write(buf[:n])
		}

		if err != nil {
			panic(err)
		}

		fmt.Println(fmt.Printf("tap2net:%d,write %d byte to network", tap2net, n))
	}
}

func (es *EdgeStar) Server(netfd io.ReadWriteCloser, tun *device.Tuntap) {

	for {
		buf := make([]byte, 2000)
		var n int
		var addr *net.UDPAddr
		var err error
		if es.Type == option.UDP {
			n, addr, err = netfd.(*net.UDPConn).ReadFromUDP(buf)
			es.Addr = addr
		} else {
			n, err = netfd.Read(buf)
		}

		if err == io.EOF {
			continue
		}

		fmt.Println(fmt.Printf("Recevied %d byte from net", n))
		//write to tap
		_, err = tun.Write(buf[:n])
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println(fmt.Printf("write %d byte to tap %s", n, tun.Name))
	}

}
