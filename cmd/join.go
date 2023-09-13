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
	"fmt"
	"github.com/spf13/cobra"
	"github.com/tiptopsoft/fvpn/pkg/device"
	"github.com/tiptopsoft/fvpn/pkg/util"
)

type joinOptions struct {
	*util.NodeCfg
	StarConfigFilePath string
	addr               string
	networkId          string
}

func joinCmd() *cobra.Command {
	var opts joinOptions
	var cmd = &cobra.Command{
		Use:          "join",
		SilenceUsage: true,
		Short:        "join a network",
		Long: `join a network which created by user, 
networkId could be found on our site after user registered, 
use free services or pay services`,

		PreRunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},

		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("networkId should be given")
			}
			return runJoin(args, &opts)
		},
	}
	fs := cmd.Flags()
	fs.StringVarP(&opts.networkId, "id", "", "", "private network id")

	return cmd
}

// runJoin join a network cmd
func runJoin(args []string, opts *joinOptions) error {
	cfg, err := util.InitConfig()
	if err != nil {
		return err
	}

	var networkId string
	if args[0] != "" {
		networkId = args[0]
	} else {
		networkId = opts.networkId
	}

	if networkId == "" {
		return errors.New("networkId is empty")
	}
	if err := device.RunJoinNetwork(cfg, networkId); err != nil {
		return err
	}

	fmt.Println(fmt.Sprintf("Join to network: %s successed", networkId))
	return nil
}
