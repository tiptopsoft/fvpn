package tun

import (
	"fmt"
	"testing"
)

func TestNew(t *testing.T) {
	dev, _ := New(2)
	fmt.Println(dev)
}
