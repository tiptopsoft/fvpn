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

package cmd

import (
	"github.com/spf13/cobra"
	"github.com/tiptopsoft/fvpn/pkg/device"
	"github.com/tiptopsoft/fvpn/pkg/util"
)

func statusCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "status",
		Short: "check fvpn status",
		Long:  "check fvpn status, if fvpn node is running will return a message, if not, will return other.",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runStatus()
		},
	}

	return cmd
}

func runStatus() error {
	cfg, err := util.InitConfig()
	if err != nil {
		return err
	}
	return device.Status(cfg.NodeCfg)
}
