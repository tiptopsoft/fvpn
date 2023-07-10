package util

import (
	"fmt"
	"testing"
)

func TestNew(t *testing.T) {
	a, _ := New(2)
	fmt.Println(a)
}
