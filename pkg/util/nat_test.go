package util

import (
	"fmt"
	"github.com/topcloudz/fvpn/pkg/packet/notify"
	"testing"
)

func TestCheckNatType(t *testing.T) {

	buff := []byte{1, 100, 0, 13, 150, 20, 31, 112, 92, 129, 204, 193, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 255, 255, 192, 168, 0, 6, 23, 173, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 255, 255, 192, 168, 0, 9}
	np, err := notify.Decode(buff)
	fmt.Println(np.DestAddr, err)
}
