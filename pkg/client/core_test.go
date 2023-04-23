package client

import (
	"fmt"
	"sync"
	"testing"
)

func TestMap(t *testing.T) {
	var m sync.Map
	m.Range(func(key, value any) bool {
		fmt.Println(value)
		return true
	})
}
