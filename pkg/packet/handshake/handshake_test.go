package handshake

import (
	"fmt"
	"strconv"
	"strings"
	"testing"
)

func TestDecode(t *testing.T) {
	s := "1 100 0 6 18 52 86 120 154 188 222 240 0 0 0 0 0 0 0 0 0 0 255 255 192 168 0 2 0 0 0 0 0 0 0 0 0 0 255 255 111 231 18 111 57 13 94 90 141 133 225 97 225 28 200 27 17 244 196 206 255 139 60 139 194 147 194 67 198 207 198 212 79 61 65 44"
	arr := strings.Split(s, " ")
	var buff []byte
	for _, v := range arr {
		value, _ := strconv.Atoi(v)
		buff = append(buff, byte(value))
	}

	pack, err := Decode(buff)
	fmt.Println(pack, err)
}
