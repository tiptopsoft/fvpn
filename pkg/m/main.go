package main

import (
	"fmt"
	"sync"
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

	var m sync.Map
	m.Store("foo", "bar")
	m.Range(func(key, value any) bool {
		fmt.Println(key, value)
		return true
	})
}
