package addr

import (
	"fmt"
	"testing"
)

func TestGenerateIP(t *testing.T) {
	e1, _ := New("aa:bb:cc:01:02:10")
	e2, _ := New("aa:bb:cc:01:02:11")
	e3, _ := New("aa:bb:cc:01:02:11")
	e4, _ := New("aa:bb:cc:01:02:11")
	e5, _ := New("aa:bb:cc:01:02:11")
	fmt.Println(e1, e2, e3, e4, e5)
}
