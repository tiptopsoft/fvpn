package processor

import (
	"context"
	"fmt"

	"github.com/interstellar-cloud/star/pkg/handler"
	"github.com/interstellar-cloud/star/pkg/log"
	"github.com/interstellar-cloud/star/pkg/option"
	"github.com/interstellar-cloud/star/pkg/processor"
	"github.com/interstellar-cloud/star/pkg/socket"
	"github.com/interstellar-cloud/star/pkg/tuntap"
	"github.com/interstellar-cloud/star/pkg/util"
)

var (
	logger = log.Log()
)

type DeviceProcessor struct {
	device *tuntap.Tuntap
	h      handler.Handler
}

func New(device *tuntap.Tuntap, h handler.Handler) processor.Processor {
	return DeviceProcessor{
		device: device,
		h:      h,
	}
}

func (dp DeviceProcessor) Process(sket socket.Interface) error {
	ctx := context.Background()
	b := make([]byte, option.STAR_PKT_BUFF_SIZE)
	size, err := dp.device.Read(b)
	destMac := util.GetMacAddr(b)
	fmt.Println(fmt.Sprintf("Read %d bytes from device %s, will write to dest %s", size, dp.device.Name, destMac))
	if err != nil {
		logger.Errorf("tap read failed. (%v)", err)
	}

	dp.h.Handle(ctx, b)
}
