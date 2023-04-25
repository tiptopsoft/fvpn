package main

import (
	"encoding/hex"
	"fmt"
)

func main() {

	str := "ff68b4ff"
	b, _ := hex.DecodeString(str)
	encodedStr := hex.EncodeToString(b)
	fmt.Printf("@@@@--bytes-->%v, length: %d\n", b, len(b))
	fmt.Printf("@@@@--string-->%s \n", encodedStr)
}
