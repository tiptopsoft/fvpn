package main

import (
	"encoding/hex"
	"fmt"
	option2 "github.com/topcloudz/fvpn/pkg/option"
	"github.com/topcloudz/fvpn/pkg/packet"
)

func main() {
	b, _ := packet.NewHeader(option2.MsgTypeRegister, "c04d6b84fd4fc978")
	buff, _ := b.Encode()
	fmt.Println(buff)

	//h, _ := b.Decode(buff)
	//r := h.(*packet.Header)
	//fmt.Printf("")

	bs, _ := hex.DecodeString("c04d6b84fd4fc978")
	fmt.Println(bs)

	fmt.Println(hex.EncodeToString(bs))

}
