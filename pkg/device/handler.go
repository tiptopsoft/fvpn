// Copyright 2023 TiptopSoft, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package device

import (
	"context"
	"github.com/tiptopsoft/fvpn/pkg/device/conn"
	"github.com/tiptopsoft/fvpn/pkg/packet/peer"
)

type Handler interface {
	Handle(ctx context.Context, frame *Frame) error
}

type HandlerFunc func(context.Context, *Frame) error

func (f HandlerFunc) Handle(ctx context.Context, frame *Frame) error {
	return f(ctx, frame)
}

type Middleware func(Handler) Handler

// Chain wrap middleware in order execute
func Chain(middlewares ...Middleware) func(Handler) Handler {
	return func(h Handler) Handler {
		for i := len(middlewares) - 1; i >= 0; i-- {
			h = middlewares[i](h)
		}

		return h
	}
}

func WithMiddlewares(handler Handler, middlewares ...Middleware) Handler {
	return Chain(middlewares...)(handler)
}

func (n *Node) handleQueryPeers(frame *Frame) {
	defer n.PutFrame(frame)
	peers, _ := peer.Decode(frame.Packet[:])
	logger.Debugf("list peers from registry: %v", peers.Peers)
	for _, info := range peers.Peers {
		dstIP := info.IP.String()
		uid := frame.UserIdString()
		if dstIP == n.device.IPToString() {
			//go over if dstIP is local dstIP
			continue
		}

		addr := info.RemoteAddr.String()
		p := n.NewPeer(uid, dstIP, n.privateKey.NewPubicKey(), n.cache)
		p.SetEndpoint(conn.NewEndpoint(addr))

		p.mode = 1
		if !p.GetStatus() {
			p.Start()
		}
	}
}
