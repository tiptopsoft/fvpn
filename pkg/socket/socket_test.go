package socket

import (
	"fmt"
	"github.com/interstellar-cloud/star/pkg/option"
	"testing"
)

func TestSocketFD(t *testing.T) {

	s := Socket{
		AppType:   "",
		Fd:        0,
		UdpSocket: nil,
	}

	fmt.Println(s.AppType == option.UDP)

	if s.AppType == option.UDP {
		fmt.Println("check ok.")
	}
}
