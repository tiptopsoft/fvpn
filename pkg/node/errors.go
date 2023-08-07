// Copyright 2023 Tiptopsoft, Inc.
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

package node

import "errors"

var (
	ErrNotImplemented = errors.New("not implement yet")
	ErrUnsupported    = errors.New("unsupported")
	ErrUnknow         = errors.New("unknown")
	ErrGetMac         = errors.New("invalid mac addr")
	ErrPacket         = errors.New("invalid packet")
	NoSuchInterface   = errors.New("route cidr+net: no such network interface")
	ErrInvalieCIDR    = errors.New("invalid cidr")
	ErrNotFound       = errors.New("not found")
)

func New(msg string) error {
	return errors.New(msg)
}
