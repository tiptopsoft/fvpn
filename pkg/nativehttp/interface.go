package nativehttp

import "github.com/topcloudz/fvpn/pkg/cache"

type Interface interface {
	ListNodes(userId string) []cache.Peer
}
