package handler

import "context"

type Interface interface {
	Handle(ctx context.Context, buff []byte) error
}
