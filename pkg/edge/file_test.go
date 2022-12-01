package edge

import (
	"fmt"
	"os"
	"testing"
)

func TestRead(t *testing.T) {
	f, _ := os.OpenFile("test", os.O_RDWR, 0)
	b := make([]byte, 1024)
	_, err := f.Read(b)
	if err != nil {
		fmt.Println(err)
	}
}
