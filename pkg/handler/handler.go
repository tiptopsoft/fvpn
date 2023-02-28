package handler

import (
	"context"
)

func Chains(ctx context.Context, buff []byte, handlers ...Interface) error {
	for _, h := range handlers {
		err := h.Handle(ctx, buff)
		if err != nil {
			return err
		}
	}

	return nil
}
