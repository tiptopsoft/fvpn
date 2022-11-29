package edge

import (
	"fmt"
	"testing"
)

func TestGetLocalMacAddr(t *testing.T) {
	fmt.Println(GetLocalMacAddr())
}
