package service

import (
	"github.com/interstellar-cloud/star/device"
	"github.com/interstellar-cloud/star/packet"
	"io"
	"net"
)

type Handler struct {
	io.ReadWriteCloser
}

func New() *Handler {
	rh := &rdHandler{}
	return &Handler{
		rh,
	}
}

// rdHandler impl read write.
type rdHandler struct {
	conn *net.UDPConn
}

func (h *rdHandler) Read(p []byte) (int, error) {
	data := make([]byte, 1024)
	n, addr, err := h.conn.ReadFromUDP(b)
	if err != nil {
		return 0, err
	}

	//
	f, err := packet.Decode(data[:24])
	if err != nil {
		switch f.Flags {
		case device.TAP_REGISTER:
			m.Store(f.SourceMac, addr)
			break
		}
	}

	p = data[24:]
	return n - 24, nil
}

func (h *rdHandler) Write(p []byte) (int, error) {

	p1 := make([]byte, 1024)
	h.conn.WriteToUDP(p1)

}

func (h *rdHandler) Close() error {
	return h.conn.Close()
}
