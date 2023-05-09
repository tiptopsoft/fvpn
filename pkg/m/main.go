package main

import (
	"fmt"
	"github.com/topcloudz/fvpn/pkg/tuntap"
	"log"
	"net"
)

func main() {
	//mac, err := net.ParseMAC("00:00:00:00:fe:80:00:00:00:00:00:00:02:00:5e:10:00:00:00:01")
	//if err != nil {
	//	panic(err)
	//}
	//
	//ip := net.IP{192, 168, 0, 1}
	//
	//reg := register.NewPacket("c04d6b84fd4fc978", mac, ip)
	//bs, err := register.Encode(reg)
	//if err != nil {
	//	panic(err)
	//}
	//
	//fmt.Println(bs)
	//fmt.Println(len(bs))
	//
	//r, _ := register.Decode(bs)
	//fmt.Println(r)

	/*buff := []byte{}

	s := "1 100 0 5 192 77 107 132 253 79 201 120 46 186 103 254 169 66 0 0 0 0 0 0 0 0 0 0 255 255 192 168 0 1 0 0 0 0 0 0 0 0 0 0 0 0 0 0"
	vs := strings.Split(s, " ")
	for _, v := range vs {
		res, _ := strconv.Atoi(v)
		buff = append(buff, uint8(res))
	}

	fmt.Println(buff)

	res, _ := register.Decode(buff)
	fmt.Println(res.SrcMac)*/

	//var m sync.Map
	//m.Store("foo", "bar")
	//m.Range(func(key, value any) bool {
	//	fmt.Println(key, value)
	//	return true
	//})

	//networkId := "c04d6b84fd4fc978"
	////mac, _ := addr.GetHostMac()
	//client := http.New("https://www.efvpn.com")
	//req := http.JoinRequest{
	//	SrcMac:    "",
	//	NetworkId: "c04d6b84fd4fc978",
	//	Ip:        "",
	//	Mask:      "",
	//}
	//resp, err := client.JoinNetwork("1", networkId, req)
	//if err != nil {
	//	fmt.Println(err)
	//}
	//
	//fmt.Println(resp)

	fa, err := net.InterfaceByName("en0")
	var ip net.IP
	addrs, err := fa.Addrs()
	for _, address := range addrs {
		// 检查ip地址判断是否回环地址
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				ip = ipnet.IP
				fmt.Println(fa.HardwareAddr, ipnet.IP, nil)
			}
		}
	}

	tun, err := tuntap.New(tuntap.TAP, ip.String(), "255.255.255.254", "1111")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(tun)

	//ifce, err := water.New(water.Config{
	//	DeviceType: water.TAP,
	//	PlatformSpecificParams: water.PlatformSpecificParams{
	//		Name:   "tap0",
	//		Driver: water.MacOSDriverTunTapOSX,
	//	},
	//})
	//if err != nil {
	//	log.Fatal(err)
	//}
	//
	//log.Printf("Interface Name: %s\n", ifce.Name())
	//
	packet := make([]byte, 2000)
	for {
		n, err := tun.Read(packet)
		if err != nil {
			log.Fatal(err)
		}
		log.Printf("Packet Received: % x\n", packet[:n])
	}
}
