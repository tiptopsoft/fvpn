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
	"log"
)

type upOptions struct {
	*util.NodeCfg
	Daemon bool
}

func nodeCmd() *cobra.Command {
	var opts upOptions
	var cmd = &cobra.Command{
		Use:          "node",
		Aliases:      []string{"n"},
		SilenceUsage: true,
		Short:        "start up a node, for private network proxy",
		Long:         `start up a node is start a private network proxy, use fvpn, you can use any device visit your private network from any place`,

		PreRunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},

		RunE: func(cmd *cobra.Command, args []string) error {
			if opts.Daemon {
				cntxt := util.GetDaemon()

				d, err := cntxt.Reborn()
				if err != nil {
					log.Fatal("Unable to run: ", err)
				}
				if d != nil {
					return nil
				}
				defer cntxt.Release()

				log.Print("fvpn started")
			}
			return runNode(&opts)
		},
	}

	fs := cmd.Flags()
	fs.BoolVarP(&opts.Daemon, "daemon", "d", false, "run daemon")

	return cmd
}

// runEdge run a client up
func runNode(opts *upOptions) error {
	config, err := util.InitConfig()
	if err != nil {
		return err
	}

	return device.Start(config)
}
