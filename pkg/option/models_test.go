package option

import (
	"context"
	"fmt"
	"testing"
)

func TestRandMac(t *testing.T) {

	mac, err := RandMac(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(mac)
}
