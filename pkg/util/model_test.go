package util

import (
	"fmt"
	"testing"
)

func TestNewLocal(t *testing.T) {
	local, err := NewLocal()
	if err != nil {
		panic(err)
	}

	f, err := local.ReadFile()
	fmt.Println(f.AppId)
}
