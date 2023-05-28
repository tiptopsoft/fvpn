package util

import (
	"fmt"
	"strconv"
	"strings"
	"testing"
)

func TestGetFrameHeader(t *testing.T) {
	s := "138 125 230 172 19 88 74 236 220 244 12 99 8 0 69 0 0 84 178 1 0 0 64 1 71 77 192 168 0 6 192 168 0 4 8 0 95 248 34 143 0 174 100 115 67 220 0 10 225 109 8 9 10 11 12 13 14 15 16 17 18 19 20 21 22 23 24 25 26 27 28 29 30 31 32 33 34 35 36 37 38 39 40 41 42 43 44 45 46 47 48 49 50 51 52 53 54 55"
	arr := strings.Split(s, " ")
	var buff []byte
	for _, v := range arr {
		value, _ := strconv.Atoi(v)
		buff = append(buff, byte(value))
	}

	h, err := GetFrameHeader(buff)
	if err != nil {
		panic(err)
	}

	fmt.Println(h)

}
