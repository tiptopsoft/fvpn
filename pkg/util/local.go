package util

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"strings"
)

type Local struct {
	file *os.File
}

type LocalConfig struct {
	Auth   string `json:"auth,omitempty"`
	AppId  string `json:"appId,omitempty"`
	UserId string `json:"userId,omitempty"`
	Debug  bool   `json:"debug,omitempty"`
}

type LocalInfo struct {
	Username string `json:"username"`
	Password string `json:"password"`
	UserId   string `json:"userId"`
}

func getLocal(mode int) (*Local, error) {
	homeDir, err := os.UserHomeDir()
	path := filepath.Join(homeDir, ".fvpn/config.json")
	_, err = os.Stat(path)
	var file *os.File
	if os.IsNotExist(err) {
		parentDir := filepath.Dir(path)
		if err := os.MkdirAll(parentDir, 0755); err != nil {
			return nil, err
		}
		file, err = os.Create(path)
	} else {
		file, err = os.OpenFile(path, mode, 0755)
	}
	return &Local{file: file}, nil
}

func (l *Local) ReadFile() (config *LocalConfig, err error) {
	decoder := json.NewDecoder(l.file)
	err = decoder.Decode(&config)
	return
}

func (l *Local) WriteFile(config *LocalConfig) error {
	encoder := json.NewEncoder(l.file)
	encoder.SetIndent("", "    ")
	return encoder.Encode(config)
}

func (l *Local) Close() error {
	return l.file.Close()
}

func GetLocalConfig() (*LocalConfig, error) {
	local, err := getLocal(os.O_RDWR)
	defer local.Close()
	if err != nil {
		return nil, err
	}
	return local.ReadFile()
}

// UpdateLocalConfig update json file
func UpdateLocalConfig(newCfg *LocalConfig) error {
	local, err := getLocal(os.O_RDWR | os.O_CREATE)
	defer local.Close()
	if err != nil {
		return err
	}
	defer local.Close()

	err = local.WriteFile(newCfg)
	if err != nil {
		return err
	}

	return nil
}

// ReplaceLocalConfig update json file
func ReplaceLocalConfig(newCfg *LocalConfig) error {
	local, err := getLocal(os.O_RDWR | os.O_TRUNC)
	defer local.Close()
	if err != nil {
		return err
	}
	defer local.Close()

	err = local.WriteFile(newCfg)
	if err != nil {
		return err
	}

	return nil
}

func GetLocalUserInfo() (info *LocalInfo, err error) {
	localCfg, err := GetLocalConfig()
	if err != nil {
		return nil, err
	}

	if localCfg.Auth == "" {
		return nil, errors.New("please login first")
	}
	info = new(LocalInfo)
	values := strings.Split(localCfg.Auth, ":")
	info.Username = values[0]
	info.Password, err = Base64Decode(values[1])
	info.UserId = localCfg.UserId
	if err != nil {
		return nil, err
	}

	return info, nil
}
