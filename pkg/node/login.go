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

package node

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/tiptopsoft/fvpn/pkg/util"
	"io"
	"os"
	"path/filepath"
)

func Login(username, password string, cfg *util.NodeCfg) error {
	client := NewClient(cfg.ControlUrl())
	req := new(LoginRequest)
	req.Username = username
	req.Password = password
	resp, err := client.Login(*req)
	if err != nil {
		return err
	}

	//登陆成功
	//logger.Debugf("login success. token:%s", resp.Token)
	//write to local
	homeDir, err := os.UserHomeDir()
	path := filepath.Join(homeDir, ".fvpn/config.json")
	_, err = os.Stat(path)
	var file *os.File
	if os.IsNotExist(err) {
		parentDir := filepath.Dir(path)
		if err := os.MkdirAll(parentDir, 0755); err != nil {
			return err
		}
		file, err = os.Create(path)
	} else {
		file, err = os.OpenFile(path, os.O_RDWR, 0755)
	}
	defer file.Close()
	var local util.LocalConfig
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&local)
	if err != nil && err != io.EOF {
		return errors.New("login failed")
	}

	appId := local.AppId
	if appId == "" {
		appId = local.AppId
	}

	b := &util.LocalConfig{
		Auth:   fmt.Sprintf("%s:%s", username, util.StringToBase64(password)),
		AppId:  local.AppId,
		UserId: resp.UserId,
	}
	return util.UpdateLocalConfig(b)
}
