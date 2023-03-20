package udp

import (
	"context"
	"errors"
	"github.com/interstellar-cloud/star/pkg/handler"
	"github.com/interstellar-cloud/star/pkg/log"
	"github.com/interstellar-cloud/star/pkg/processor"
	"github.com/interstellar-cloud/star/pkg/socket"
	"github.com/interstellar-cloud/star/pkg/tuntap"
	"io"
)

var (
	logger = log.Log()
)

type UdpProcessor struct {
	h      handler.Handler
	device *tuntap.Tuntap
	skt    socket.Interface
}

func New(device *tuntap.Tuntap, h handler.Handler, skt socket.Interface) processor.Processor {
	return UdpProcessor{
		device: device,
		h:      h,
		skt:    skt,
	}
}

func (up UdpProcessor) Process() error {
	udpBytes := make([]byte, 2048)
	size, err := up.skt.Read(udpBytes)
	if size < 0 {
		return errors.New("no data exists")
	}
	logger.Infof("star net skt receive size: %d, data: (%v)", size, udpBytes[:size])
	if err != nil {
		if err == io.EOF {
			//no data exists, continue read next frame continue
			logger.Errorf("not data exists")
		} else {
			logger.Errorf("read from remote error: %v", err)
		}
	}

	return up.h.Handle(context.Background(), udpBytes[:size])
}
