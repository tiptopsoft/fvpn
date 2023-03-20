package udp

import (
	"context"
	"errors"
	"fmt"
	"io"
	"unsafe"

	"github.com/interstellar-cloud/star/pkg/handler"
	"github.com/interstellar-cloud/star/pkg/log"
	"github.com/interstellar-cloud/star/pkg/node"
	"github.com/interstellar-cloud/star/pkg/option"
	"github.com/interstellar-cloud/star/pkg/packet"
	"github.com/interstellar-cloud/star/pkg/packet/forward"
	"github.com/interstellar-cloud/star/pkg/packet/register/ack"
	"github.com/interstellar-cloud/star/pkg/processor"
	"github.com/interstellar-cloud/star/pkg/socket"
	"github.com/interstellar-cloud/star/pkg/tuntap"
	"github.com/interstellar-cloud/star/pkg/util"
)

var (
	logger = log.Log()
)

type UdpProcessor struct {
	h      handler.Handler
	device *tuntap.Tuntap
}

func New(device *tuntap.Tuntap, handler.Handler) processor.Processor {
	return UdpProcessor{
		device: device,
		h: h,
	}
}

func (up UdpProcessor) Process(skt socket.Interface) error {

	device := up.device
	if s.Protocol == option.UDP {
		udpBytes := make([]byte, 2048)
		size, err := skt.Read(udpBytes)
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

		return up.h.Handle(context.Background(), udpBytes)
	}

	return nil
}	
		