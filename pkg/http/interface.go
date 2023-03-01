package http

import "github.com/interstellar-cloud/star/pkg/node"

type Interface interface {
	ListNodes(userId string) []node.Node
}
