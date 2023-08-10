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

package device

type RouterManager interface {
	AddRouter(cidr string) error
	RemoveRouter(cidr string) error
}

type router struct {
	cidr string
	name string
}

func NewRouter(cidr, name string) RouterManager {
	return &router{
		cidr: cidr,
		name: name,
	}
}
