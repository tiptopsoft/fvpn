package http

import "github.com/interstellar-cloud/star/pkg/cache"

type Interface interface {
	ListNodes(userId string) []cache.Peer
}
