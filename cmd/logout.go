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
	"fmt"
	"github.com/spf13/cobra"
	"github.com/tiptopsoft/fvpn/pkg/util"
)

func logout() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "logout",
		Long: `logout out fvpn, once you logout, packet will be invalid, because fvpn check each packet`,
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			return runLogout()
		},
	}

	return cmd
}

func runLogout() error {
	localCfg, err := util.GetLocalConfig()
	if err != nil {
		return fmt.Errorf("logout failed, %v", err)
	}
	localCfg.Auth = ""
	localCfg.UserId = ""
	return util.ReplaceLocalConfig(localCfg)
}
