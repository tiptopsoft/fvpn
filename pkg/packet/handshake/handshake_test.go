package handshake

import (
	"fmt"
	"strconv"
	"strings"
	"testing"
)

func TestDecode(t *testing.T) {
	s := "1 100 0 12 0 0 0 0 0 0 0 0 18 52 86 120 154 188 222 240 0 0 0 0 0 0 0 0 0 0 255 255 5 244 24 141 227 55 129 157 0 239 22 167 212 226 196 118 40 16 105 105 143 100 2 53 132 238 79 128 16 80 28 120 89 215 225 50 0 0 0 0 0 0 0 0 0 0 0 0"
	arr := strings.Split(s, " ")
	var buff []byte
	for _, v := range arr {
		value, _ := strconv.Atoi(v)
		buff = append(buff, byte(value))
	}

	pack, err := Decode(buff)
	fmt.Println(pack, err)
}
