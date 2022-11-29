package registry

import (
	"fmt"
	"strconv"
	"testing"
)

func TestGenerateIP(t *testing.T) {
	ip1 := string2Long("192.168.0.1")
	fmt.Println(ip1)
	fmt.Println(strconv.FormatInt(36, 10))
	fmt.Println(GenerateIP(ip1))
}
