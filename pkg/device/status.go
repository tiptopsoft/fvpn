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
	"fmt"
	"github.com/olekukonko/tablewriter"
	"github.com/tiptopsoft/fvpn/pkg/util"
	"os"
)

func Status(cfg *util.NodeCfg) error {
	client := NewClient(cfg.HostUrl())
	resp, err := client.Status()
	if err != nil {
		return err
	}

	var data [][]string
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Status", "Version"})
	if resp != nil {
		data = append(data, []string{fmt.Sprintf("%s", resp.Status), resp.Version})
	}
	table.AppendBulk(data)
	table.Render()
	return nil
}

func Stop(cfg *util.NodeCfg) error {
	client := NewClient(cfg.HostUrl())
	_, err := client.Stop()
	if err != nil {
		return err
	}

	return nil
}
