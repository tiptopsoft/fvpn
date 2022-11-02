package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"reflect"
)

func main() {
	//var v2 int32
	//var b [4]byte
	//
	//v2 = 257
	//
	//b[3] = uint8(v2)
	//b[2] = uint8(v2 >> 8)
	//b[1] = uint8(v2 >> 16)
	//b[0] = uint8(v2 >> 24)
	//
	//fmt.Println(b)
	//
	//b2 := IntToBytes(257)
	//fmt.Println("b2:", b2)

	s := "魑"
	//for i, j := range s {
	//	fmt.Println(reflect.TypeOf(s[i]))
	//	fmt.Println(i)
	//	fmt.Println(reflect.TypeOf(j))
	//}
	sb := []byte(s)
	fmt.Println(sb)
	sbb := []rune(s)
	fmt.Println(sbb)
	fmt.Println(reflect.TypeOf(s[0]))

	fmt.Println(sbb[0])
	fmt.Print(string([]rune{39761}))
}

func IntToBytes(n int) []byte {
	x := int32(n)
	bytesBuffer := bytes.NewBuffer([]byte{})
	binary.Write(bytesBuffer, binary.BigEndian, x)
	return bytesBuffer.Bytes()
}
