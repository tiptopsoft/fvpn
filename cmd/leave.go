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

package cmd

import (
	"errors"
	"github.com/spf13/cobra"
	"github.com/tiptopsoft/fvpn/pkg/device"
	"github.com/tiptopsoft/fvpn/pkg/util"
)

// Join a networkId will be created tun device. and will be assigned a IP which can be found in our website.
type leaveOptions struct {
	*util.NodeCfg
	StarConfigFilePath string
	NetworkId          string
}

func leaveCmd() *cobra.Command {
	var opts joinOptions
	var cmd = &cobra.Command{
		Use:          "leave",
		SilenceUsage: true,
		Short:        "leave a network",
		Long: `leave a joined network, once use leave a network, 
fvpn can not route any frame to dst node, 
if you want continue your destination routing, 
you can join it again`,

		PreRunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},

		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("networkId should be given")
			}
			return runLeave(args)
		},
	}
	fs := cmd.Flags()
	fs.StringVarP(&opts.networkId, "id", "", "", "network id")

	return cmd
}

// runJoin join a network cmd
func runLeave(args []string) error {

	cfg, err := util.InitConfig()
	if err != nil {
		return err
	}
	return device.RunLeaveNetwork(cfg, args[0])
}
