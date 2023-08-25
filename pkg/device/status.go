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
	"fmt"
	"github.com/tiptopsoft/fvpn/pkg/util"
)

func Status(cfg *util.NodeCfg) error {
	client := NewClient(cfg.HostUrl())
	resp, err := client.Status()
	if err != nil {
		return err
	}

	if resp == nil || resp.Status == "" {
		fmt.Println("fvpn not running, please check")
	} else {
		fmt.Println(fmt.Sprintf("Status: %s, Version: %s", resp.Status, resp.Version))
	}
	return nil
}

func Stop(cfg *util.NodeCfg) error {
	client := NewClient(cfg.HostUrl())
	resp, err := client.Stop()
	if err != nil {
		return err
	}

	fmt.Println(resp.Result)
	return nil
}
