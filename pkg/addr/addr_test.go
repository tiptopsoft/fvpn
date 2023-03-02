package addr

import (
	"fmt"
	"net"
	"testing"
)

func TestNew(t *testing.T) {
	RecMac := "01:01:03:02:03:01"
	src, err := net.ParseMAC(RecMac)
	if err != nil {
		t.Errorf("%v", err)
	}
	endpoint, _ := New(src)
	fmt.Println(endpoint)
}
