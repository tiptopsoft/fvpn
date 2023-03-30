package option

import (
	"bytes"
	"errors"
	"github.com/spf13/viper"
)

type Protocol string

const (
	TCP Protocol = "tcp"
	UDP Protocol = "udp"
)

var (
	STAR_PKT_BUFF_SIZE = 2048
	defaultYaml        = []byte(`star:
  listen: :3000
  registry: :4000
  tap: tap0
  ip: 192.168.0.1
  mask: 255.255.255.0
  mac: 01:02:0f:0E:04:01
  type: udp

#-------------------分割线
registry:
  listen: 127.0.0.1:4000
  httpListen: :4009
  type: udp`)
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
	Star *StarConfig `mapstructure:"star"`
	Reg  *RegConfig  `mapstructure:"registry"`
}

// StarConfig read from a config file or cmd flags, or can be assgined from a registry after got the registry ack.
type StarConfig struct {
	Registry    string   `mapstructure:"registry"`
	Listen      string   `mapstructure:"listen"`
	TapName     string   `mapstructure:"tap"`
	TapIP       string   `mapstructure:"ip"`
	TapMask     string   `mapstructure:"mask"`
	MacAddr     string   `mapstructure:"mac"`
	Protocol    Protocol `mapstructure:"type"`
	OpenAuth    bool
	OpenEncrypt bool
}

type Mysql struct {
	User     string `mapstructure:"user"`
	Url      string `mapstructure:"url"`
	Password string `mapstructure:"password"`
	Name     string `mapstructure:"name"`
}

type RegConfig struct {
	Listen      string   `mapstructure:"listen"`
	HttpListen  string   `mapstructure:"httpListen"`
	Protocol    Protocol `mapstructure:"type"`
	OpenAuth    bool
	OpenEncrypt bool
}

func InitConfig() (config *Config, err error) {
	viper.SetConfigName("app")         // name of config file (without extension)
	viper.SetConfigType("yaml")        // REQUIRED if the config file does not have the extension in the name
	viper.AddConfigPath("/etc/star/")  // path to look for the config file in
	viper.AddConfigPath("$HOME/.star") // call multiple times to add many search paths
	viper.AddConfigPath(".")           // optionally look for config in the working directory
	viper.AddConfigPath("./conf/")
	if err = viper.ReadInConfig(); err != nil { // Handle errors reading the config file
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			viper.ReadConfig(bytes.NewBuffer(defaultYaml))
		} else {
			return nil, errors.New("invalid config")
		}
	}
	if err = viper.UnmarshalExact(&config); err != nil {
		return nil, err
	}

	return

}
