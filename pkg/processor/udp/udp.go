package udp

import (
	"github.com/interstellar-cloud/star/pkg/handler"
	"github.com/interstellar-cloud/star/pkg/processor"
	"github.com/interstellar-cloud/star/pkg/socket"
)

type UdpProcessor struct {
	h handler.Handler
}

func New(h handler.Handler) processor.Processor {
	return UdpProcessor{
		h: h,
	}
}

func (up UdpProcessor) Process(sk socket.Interface) {

}
