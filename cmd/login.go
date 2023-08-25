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
	"bufio"
	"errors"
	"fmt"
	"github.com/moby/term"
	"github.com/spf13/cobra"
	"github.com/tiptopsoft/fvpn/pkg/device"
	"github.com/tiptopsoft/fvpn/pkg/util"
	"os"
)

// Join a networkId will be created tun device. and will be assigned a IP which can be found in our website.
type loginOptions struct {
	Username string
	Password string
}

func loginCmd() *cobra.Command {
	var opts loginOptions
	var cmd = &cobra.Command{
		Use:          "login",
		SilenceUsage: true,
		Short:        "login fvpn",
		Long:         `when you are using fvpn, you should logon first, login fvpn use username and password which registered on our site, if you did not logon, you can not join any networks.`,

		PreRunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},

		RunE: func(cmd *cobra.Command, args []string) error {
			return runLogin(opts)
		},
	}
	fs := cmd.Flags()
	fs.StringVarP(&opts.Username, "username", "u", "", "username for fvpn")
	fs.StringVarP(&opts.Password, "password", "p", "", "username for fvpn")

	return cmd
}

// runJoin join a network cmd
func runLogin(opts loginOptions) error {
	config, err := util.InitConfig()
	if err != nil {
		return err
	}

	//s := &client.Peer{
	//	Config: config,
	//}

	if opts.Password == "" {
		if opts.Username == "" {
			opts.Username, _ = readLine("Username: ", false)
		}

		if opts.Username == "" {
			if token, err := readLine("Token", true); err != nil {
				return errors.New("token required")
			} else {
				opts.Password = token
			}
		} else {
			if password, err := readLine("Password: ", true); err != nil {
				return errors.New("password required")
			} else {
				opts.Password = password
			}
		}
	}

	//check whether has been login TODO

	err = device.Login(opts.Username, opts.Password, config.NodeCfg)
	if err != nil {
		return err
	}
	fmt.Println("Login Succeeded")

	return nil

}

//func readLine(prompt string, slient bool) (string, error) {
//	fmt.Print(prompt)
//	reader := bufio.NewReader(os.Stdin)
//	line, _, err := reader.ReadLine()
//	if err != nil {
//		panic(err)
//		return "", err
//	}
//
//	return string(line), err
//}

func readLine(prompt string, slient bool) (string, error) {
	fmt.Print(prompt)
	if slient {
		fd := os.Stdin.Fd()
		state, err := term.SaveState(fd)
		if err != nil {
			return "", err
		}
		term.DisableEcho(fd, state)
		defer term.RestoreTerminal(fd, state)
	}

	reader := bufio.NewReader(os.Stdin)
	line, _, err := reader.ReadLine()
	if err != nil {
		return "", err
	}
	if slient {
		fmt.Println()
	}

	return string(line), nil
}
