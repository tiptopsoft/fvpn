package util

import (
	"encoding/json"
	"os"
	"path/filepath"
)

type Local struct {
	file *os.File
}

type LocalConfig struct {
	Auth  string `json:"auth"`
	AppId string `json:"appId"`
}

type Response struct {
	Code    int         `json:"code"`
	Result  interface{} `json:"data"`
	Message string      `json:"msg"`
}

type JoinRequest struct {
	NetWorkId string `json:"networkId"`
	CIDR      string `json:"cidr"`
}

type JoinResponse struct {
	CIDR    string `json:"cidr"`
	IP      string `json:"ip"`
	Name    string `json:"name"`
	Network string `json:"network"`
}

type InitResponse struct {
	IP    string `json:"ip"`
	Mask  string `json:"mask"`
	AppId string `json:"appId"`
}

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type LoginResponse struct {
	Token string `json:"token"`
}

func HttpOK(data interface{}) Response {
	return Response{
		Code:    200,
		Result:  data,
		Message: "",
	}
}

func HttpError(message string) Response {
	return Response{
		Code:    500,
		Result:  nil,
		Message: message,
	}
}
func NewLocal(mode int) (*Local, error) {
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
	return encoder.Encode(config)
}

func (l *Local) Close() error {
	return l.file.Close()
}

//func GetLocalConfig() (*LocalConfig, error) {
//	homeDir, err := os.UserHomeDir()
//	path := filepath.Join(homeDir, ".fvpn/config.json")
//	_, err = os.Stat(path)
//	var file *os.file
//	if os.IsNotExist(err) {
//		parentDir := filepath.Dir(path)
//		if err := os.MkdirAll(parentDir, 0755); err != nil {
//			return nil, err
//		}
//		file, err = os.Create(path)
//	} else {
//		file, err = os.OpenFile(path, os.O_RDWR, 0755)
//	}
//	defer file.Close()
//	var local LocalConfig
//	decoder := json.NewDecoder(file)
//	err = decoder.Decode(&local)
//	if err != nil {
//		return nil, errors.New("login failed")
//	}
//
//	return &local, nil
//}
