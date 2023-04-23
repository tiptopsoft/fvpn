package executor

import (
	"github.com/topcloudz/fvpn/pkg/socket"
)

type Executor interface {
	Execute(socket socket.Interface) error
}
