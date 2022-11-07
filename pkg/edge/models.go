package edge

import (
	"context"
	"fmt"
	"github.com/interstellar-cloud/star/pkg/internal"
	"github.com/spf13/viper"
	"math/rand"
)

// EdgeConfig read from a config file or cmd flags, or can be assgined from a registry after got the register ack.
type EdgeConfig struct {
	Registry string `mapstructure:"registies"`
	Listen   string `mapstructure:"listen"`
	TapName  string
	TapIP    string
	TapMask  string
	MacAddr  string
	Protocol int
}

func EdgeDefault() *EdgeConfig {
	return &EdgeConfig{}
}

const (
	TCP = iota
	UDP
)

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
	Listen string `mapstructure:"listen"`
	Mysql  Mysql  `mapstructure:"mysql"`
	Proto  internal.Protocol
}

type Mysql struct {
	User     string `mapstructure:"user"`
	Url      string `mapstructure:"url"`
	Password string `mapstructure:"password"`
	Name     string `mapstructure:"name"`
}

// RandMac rand gen a mac
func RandMac(ctx context.Context) (string, error) {
	buf := make([]byte, 6)
	_, err := rand.Read(buf)
	if err != nil {
		return "", err
	}
	buf[0] |= 2
	mac := fmt.Sprintf("%02x:%02x:%02x:%02x:%02x:%02x", buf[0], buf[1], buf[2], buf[3], buf[4], buf[5])

	return mac, nil
}

func InitConfig() (config *Config, err error) {
	viper.SetConfigName("app")                  // name of config file (without extension)
	viper.SetConfigType("yaml")                 // REQUIRED if the config file does not have the extension in the name
	viper.AddConfigPath("/etc/edge/")           // path to look for the config file in
	viper.AddConfigPath("$HOME/.edge")          // call multiple times to add many search paths
	viper.AddConfigPath(".")                    // optionally look for config in the working directory
	if err = viper.ReadInConfig(); err != nil { // Handle errors reading the config file
		panic(fmt.Errorf("fatal error config file: %w", err))
	}

	if err = viper.UnmarshalExact(&config); err != nil {
		return nil, err
	}

	return

}
