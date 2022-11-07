package edge

import (
	"context"
	"fmt"
	"github.com/interstellar-cloud/star/pkg/option"
	"testing"
)

func TestRandMac(t *testing.T) {

	mac, err := option.RandMac(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(mac)
}
