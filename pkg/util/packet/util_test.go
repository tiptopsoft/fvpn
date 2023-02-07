package packet

import (
	"fmt"
	"github.com/interstellar-cloud/star/pkg/util"
	"testing"
)

func TestGetLocalMac(t *testing.T) {
	ac, _ := util.GetLocalMac("en0")
	fmt.Println(ac)
}
