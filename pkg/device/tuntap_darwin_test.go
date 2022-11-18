package device

import (
	"fmt"
	"testing"
)

func TestNew(t *testing.T) {
	t1, err := New(TAP)
	if err != nil {
		panic(err)
	}

	fmt.Println(t1)
}
