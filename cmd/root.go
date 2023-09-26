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
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:          "fvpn [command]",
	SilenceUsage: true,
	Short:        "fvpn can let you join/leave a private network, compose network",
	Long:         `fvpn can let you join/leave a private network, compose our network, build node mesh and so on`,
}

func Execute() {
	rootCmd.AddCommand(nodeCmd(), registryCmd(), joinCmd(), loginCmd(), logout(), leaveCmd(), statusCmd(), stopCmd(), listCmd())
	if err := rootCmd.Execute(); err != nil {
		os.Exit(-1)
	}
}
