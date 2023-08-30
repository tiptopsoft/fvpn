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

package util

import (
	"fmt"
	"github.com/spf13/viper"
)

type Protocol string

type SuperStar struct {
	Listen int
}

type Star struct {
	Name string
	IP   string
	Mask string
	Mode int //0 tun 1 tap
}

type Config struct {
	NodeCfg     *NodeCfg     `mapstructure:"node"`
	RegistryCfg *RegistryCfg `mapstructure:"registry"`
}

// NodeCfg read from a config file or cmd flags, node and registry configurations are here.
type NodeCfg struct {
	Registry   string `mapstructure:"registry"`
	Udp        string `mapstructure:"udp"`
	Listen     int    `mapstructure:"listen"`
	HttpListen string `mapstructure:"httpListen"`
	ConsoleUrl string `mapstructure:"consoleUrl"`
	//Protocol   Protocol `mapstructure:"type"`
	Offset  int32   `mapstructure:"offset"`
	Encrypt Encrypt `mapstructure:"encrypt"`
	Auth    Auth    `mapstructure:"auth"`
	Relay   Relay   `mapstructure:"relay"`
	Log     LogCfg  `mapstructure:"log"`
	IPV6    IPV6Cfg `mapstructure:"ipv6"`
	Driver  string  `mapstructure:"driver"`
}

type IPV6Cfg struct {
	Enable bool `mapstructure:"enable"`
}

type LogCfg struct {
	EnableDebug bool `mapstructure:"debug"`
}

type Encrypt struct {
	Enable bool `mapstructure:"enable"`
}

type Auth struct {
	Enable bool `mapstructure:"enable"`
}

type Relay struct {
	Enable bool `mapstructure:"enable"`
}

func (cfg *NodeCfg) EnableRelay() bool {
	return cfg.Relay.Enable
}

func (cfg *NodeCfg) EnableAuth() bool {
	return cfg.Auth.Enable
}

func (cfg *NodeCfg) EnableEncrypt() bool {
	return cfg.Encrypt.Enable
}

func (cfg *NodeCfg) HttpListenStr() string {
	return cfg.HttpListen
}

func (cfg *NodeCfg) HostUrl() string {
	return fmt.Sprintf("http://127.0.0.1%s", cfg.HttpListen)
}

func (cfg *NodeCfg) ControlUrl() string {
	return cfg.ConsoleUrl
}

func (cfg *NodeCfg) RegistryUrl() string {
	return cfg.Registry
}

func (cfg *NodeCfg) AuthEnable() bool {
	return cfg.Auth.Enable
}

type Redis struct {
	Enable   bool   `json:"enable"`
	Url      string `json:"host"`
	Port     int    `json:"port"`
	Username string `json:"username"`
	Password string `json:"password"`
}

type RegistryCfg struct {
	Listen     string `mapstructure:"Listen"`
	HttpListen string `mapstructure:"HttpListen"`
	Redis      Redis  `json:redis`
	Driver     string `mapstructure:"driver"`
}

func InitConfig() (config *Config, err error) {
	viper.SetConfigName("app")         // name of config file (without extension)
	viper.SetConfigType("yaml")        // REQUIRED if the config file does not have the extension in the name
	viper.AddConfigPath("/etc/fvpn/")  // path to look for the config file in
	viper.AddConfigPath("$HOME/.fvpn") // call multiple times to add many search paths
	viper.AddConfigPath(".")           // optionally look for config in the working directory
	viper.AddConfigPath("./conf/")

	// read default
	//viper.ReadConfig(bytes.NewBuffer(defaultYaml))
	viper.SetDefault("node.ConsoleUrl", "https://www.tiptopsoft.cn")
	viper.SetDefault("node.HttpListen", ":6662")
	viper.SetDefault("node.Registry", "tiptopsoft.cn")
	viper.SetDefault("node.Relay.Enable", true)
	viper.SetDefault("node.log.debug", false)
	viper.SetDefault("node.IPV6.Enable", false)
	viper.SetDefault("node.Udp", "udp4")
	viper.SetDefault("node.Encrypt.Enable", true)
	viper.SetDefault("node.Auth.Enable", true)
	viper.SetDefault("node.Listen", 6061)
	viper.SetDefault("registry.Listen", ":4000")
	viper.SetDefault("registry.HttpListen", ":4001")
	viper.SetDefault("registry.redis.enable", false)
	if err = viper.ReadInConfig(); err != nil { // Handle errors reading the config file
		//if _, ok := err.(viper.ConfigFileNotFoundError); ok {
		//	viper.ReadConfig(bytes.NewBuffer(defaultYaml))
		//} else {
		//	return nil, errors.New("invalid config")
		//}

	}

	if err = viper.UnmarshalExact(&config); err != nil {
		return nil, err
	}

	return

}
