package processor

import (
	"context"
	"fmt"

	"github.com/topcloudz/fvpn/pkg/handler"
	"github.com/topcloudz/fvpn/pkg/log"
	"github.com/topcloudz/fvpn/pkg/option"
	"github.com/topcloudz/fvpn/pkg/processor"
	"github.com/topcloudz/fvpn/pkg/tuntap"
	"github.com/topcloudz/fvpn/pkg/util"
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

func (dp DeviceProcessor) Process() error {
	ctx := context.Background()
	b := make([]byte, option.FVPN_PKT_BUFF_SIZE)
	size, err := dp.device.Read(b)
	destMac := util.GetMacAddr(b)
	fmt.Println(fmt.Sprintf("Read %d bytes from device %s, will write to dest %s", size, dp.device.Name, destMac))
	if err != nil {
		logger.Errorf("tap read failed. (%v)", err)
	}

	return dp.h.Handle(ctx, b[:size])
}
