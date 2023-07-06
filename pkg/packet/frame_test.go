package packet

import (
	"fmt"
	"testing"
)

func TestNewFrame(t *testing.T) {
	f := NewFrame()
	s := "abc"
	buff := []byte(s)
	size := len(buff)
	copy(f.Packet, buff)

	fmt.Println("packet: ", f.Packet)

	b := "cde"
	buff1 := []byte(b)

	copy(f.Packet[size:], buff1)
	fmt.Println("packet1: ", f.Packet)
}
