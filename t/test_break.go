package main

import (
	"fmt"
)

func main() {

	i := 0
	//outloop:
	for {
		switch i {
		case 0:
			fmt.Println("case 0: ", i)
			i++
			break
		case 1:
			fmt.Println("case 1: ", i)
			i++
			break
		case 2:
			fmt.Println("case 2: ", i)
			i++
			break
		case 3:
			fmt.Println("break out for")
			break
		default:
			fmt.Println("default ...")
		}

	}
}
