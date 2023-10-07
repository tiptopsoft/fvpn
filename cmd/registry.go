// Copyright 2023 TiptopSoft, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//	http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/tiptopsoft/fvpn/pkg/relay"
	"github.com/tiptopsoft/fvpn/pkg/util"
)

type RegOptions struct {
	Listen int
}

func registryCmd() *cobra.Command {
	var opts RegOptions
	cmd := &cobra.Command{
		Use:   "serve",
		Short: "s",
		Long:  `fvpn start a registry, a data center/relay server, is our core service`,
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			return runSuper(&opts)
		},
	}

	fs := cmd.Flags()
	fs.IntVarP(&opts.Listen, "port", "p", 0, "registry server port")

	return cmd
}

func runSuper(opts *RegOptions) error {
	config, err := util.InitConfig()
	if err != nil {
		return err
	}
	s := relay.RegServer{
		RegistryCfg: config.RegistryCfg,
	}

	if opts.Listen != 0 {
		s.RegistryCfg.Listen = fmt.Sprintf(":%d", opts.Listen)
	}

	return s.Start()
}
