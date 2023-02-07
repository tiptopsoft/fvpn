package main

import (
	"fmt"
	"github.com/interstellar-cloud/star/pkg/registry/addr"
)

func main() {
	e1, _ := addr.New("aa:bb:cc:01:02:10")
	e2, _ := addr.New("aa:bb:cc:01:02:11")

	fmt.Println(e1, e2)

}
