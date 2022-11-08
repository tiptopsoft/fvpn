package option

import (
	"fmt"
	"github.com/spf13/viper"
)

type Protocol string

const (
	TCP Protocol = "tcp"
	UDP Protocol = "udp"
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
	Star *EdgeConfig `mapstructure:"star"`
	Reg  *RegConfig  `mapstructure:"register"`
}

type Mysql struct {
	User     string `mapstructure:"user"`
	Url      string `mapstructure:"url"`
	Password string `mapstructure:"password"`
	Name     string `mapstructure:"name"`
}

func InitConfig() (config *Config, err error) {
	viper.SetConfigName("app")         // name of config file (without extension)
	viper.SetConfigType("yaml")        // REQUIRED if the config file does not have the extension in the name
	viper.AddConfigPath("/etc/edge/")  // path to look for the config file in
	viper.AddConfigPath("$HOME/.edge") // call multiple times to add many search paths
	viper.AddConfigPath(".")           // optionally look for config in the working directory
	viper.AddConfigPath("./conf/")
	if err = viper.ReadInConfig(); err != nil { // Handle errors reading the config file
		panic(fmt.Errorf("fatal error config file: %w", err))
	}

	if err = viper.UnmarshalExact(&config); err != nil {
		return nil, err
	}

	return

}
