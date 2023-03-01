package executor

import (
	"github.com/interstellar-cloud/star/pkg/socket"
)

type Executor interface {
	Execute(socket socket.Interface) error
}
