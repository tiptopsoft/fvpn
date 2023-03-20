package processor

import (
	"github.com/interstellar-cloud/star/pkg/socket"
)

type Processor interface {
	Process(socket.Interface)
}
