package peer

import (
	"fmt"
	"github.com/topcloudz/fvpn/pkg/packet/peer/ack"
	"testing"
	"unsafe"
)

func TestEncode(t *testing.T) {
	fmt.Print(unsafe.Sizeof(ack.EdgeInfo{}))
}
