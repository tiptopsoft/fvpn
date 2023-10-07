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

func RunListNetworks(cfg *util.Config) error {
	logger.Debugf("start list networks")

	cm := NewManager(cfg.NodeCfg)
	resp, err := cm.ListNetworks()
	if err != nil {
		return err
	}

	var data [][]string
	if resp.List != nil {
		for index, id := range resp.List {
			data = append(data, []string{fmt.Sprintf("%d", index+1), id.NetworkId})
		}
	}

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Index", "Name"})

	table.AppendBulk(data)

	table.Render()
	return nil
}
