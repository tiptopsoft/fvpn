package executor

import (
	"github.com/interstellar-cloud/star/pkg/util/socket"
)

type Executor interface {
	Execute(socket socket.Socket) error
}
