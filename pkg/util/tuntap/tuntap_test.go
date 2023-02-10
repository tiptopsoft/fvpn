package tuntap

import (
	"fmt"
	"testing"
)

func TestNew(t *testing.T) {

	tap, err := New(TAP)
	if err != nil {
		panic(err)
	}

	for {
		b := make([]byte, 1024)
		n, err := tap.Socket.Read(b)
		if err != nil {
			fmt.Println(err)
		}

		fmt.Println(fmt.Sprintf("Read %d bytes from device %s", n, tap.Name))
	}
}
