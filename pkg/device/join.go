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

import (
	"github.com/tiptopsoft/fvpn/pkg/util"
)

func RunJoinNetwork(cfg *util.Config, networkId string) error {
	logger.Debugf("start join to network: %s", networkId)

	cm := NewManager(cfg.NodeCfg)
	resp, err := cm.JoinNetwork(networkId)
	if err != nil {
		return err
	}

	return NewRouter(resp.CIDR, resp.Name).AddRouter(resp.CIDR)
}

func RunLeaveNetwork(cfg *util.Config, networkId string) error {

	logger.Infof("start leave network: %s", networkId)

	cm := NewManager(cfg.NodeCfg)
	resp, err := cm.LeaveNetwork(networkId)
	if err != nil {
		return err
	}

	return NewRouter(resp.CIDR, resp.Name).RemoveRouter(resp.CIDR)
}
