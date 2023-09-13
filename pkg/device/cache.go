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

type Interface interface {
	Set(userId, ip string, peer *Peer) error
	Get(userId, ip string) (*Peer, error)
	Lists(userId string) PeerMap
}

func NewCache(driver string) Interface {
	if driver == "" {
		driver = "local"
	}

	switch driver {
	case "local":
		return newLocal()
	}

	return nil
}
