package packet

import (
	"github.com/interstellar-cloud/star/pkg/log"
	"net"
	"sync"
)

var socket sync.Map
var once sync.Once

type StarSock struct {
	Family uint8
	Type   uint8
	Port   uint16
	Addr   [128]byte
}

func SendPacket(data []byte, dst string) error {
	if sock, ok := socket.Load(dst); !ok {
		once.Do(func() {
			conn, err := net.Dial("udp", dst)
			if err != nil {
				log.Logger.Error("connect dst failed.err: %v", err)
			}

			socket.Store(dst, conn)
			sock = conn
		})

		if _, err := sock.(*net.UDPConn).Write(data); err != nil {
			log.Logger.Error("write dst failed.err: %v", err)
		}

	}

	return nil
}
