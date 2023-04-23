package udp

import (
	"context"
	"errors"
	"github.com/topcloudz/fvpn/pkg/handler"
	"github.com/topcloudz/fvpn/pkg/log"
	"github.com/topcloudz/fvpn/pkg/processor"
	"github.com/topcloudz/fvpn/pkg/socket"
	"io"
)

var (
	logger = log.Log()
)

type Processor struct {
	h   handler.Handler
	skt socket.Interface
}

func New(h handler.Handler, skt socket.Interface) processor.Processor {
	return Processor{
		h:   h,
		skt: skt,
	}
}

func (up Processor) Process() error {
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
