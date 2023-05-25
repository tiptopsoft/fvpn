package util

import (
	"fmt"
	"github.com/topcloudz/fvpn/pkg/packet/notify"
	"strconv"
	"strings"
	"testing"
)

func TestCheckNatType(t *testing.T) {

	/*np := notify.NewPacket("96141f705c81ccc1")
	np.SourceIP = net.ParseIP("192.168.0.9")
	np.Port = 6061
	np.NatIP = np.SourceIP
	np.NatType = option.RestrictNat
	np.NatPort = 34343
	np.DestAddr = net.ParseIP("192.168.0.6")
	buff1, _ := notify.Encode(np)
	fmt.Println(buff1)
	//s := "1 100 0 13 150 20 31 112 92 129 204 193 0 0 0 0 0 0 0 0 0 0 255 255 192 168 0 9 23 173 0 0 1 0 0 0 0 0 0 0 0 0 0 255 255 192 168 0 6 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0"
	//arr := strings.Split(s, " ")
	//var buff []byte
	//for _, v := range arr {
	//	value, _ := strconv.Atoi(v)
	//	buff = append(buff, byte(value))
	//}

	//fmt.Println(len(buff))
	//buff := []byte{1, 100, 0, 13, 150, 20, 31, 112, 92, 129, 204, 193, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 255, 255, 192, 168, 0, 6, 23, 173, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 255, 255, 192, 168, 0, 9}
	np, err := notify.Decode(buff1)
	fmt.Println(np, err)*/

	s := "1 100 0 13 150 20 31 112 92 129 204 193 0 0 0 0 0 0 0 0 0 0 255 255 192 168 0 9 0 0 1 0 0 0 0 0 0 0 0 0 0 255 255 81 70 36 156 0 0 0 0 0 0 0 0 0 0 255 255 81 70 36 156 23 173 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0"
	arr := strings.Split(s, " ")
	var buff []byte
	for _, v := range arr {
		value, _ := strconv.Atoi(v)
		buff = append(buff, byte(value))
	}

	np, err := notify.Decode(buff)
	fmt.Println(np, err)
}
