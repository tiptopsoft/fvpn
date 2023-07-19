package util

import (
	"fmt"
	"github.com/spf13/viper"
)

type Protocol string

var (
	defaultYaml = []byte(
		`client:
  Listen: :3000
  HttpListen: :6662
#  server: 127.0.0.1:4000
  server: 211.159.225.186:4000
  ConsoleUrl: https://www.efvpn.com
  type: udp
  offset: 1
  Relay:
    Enable: true
  Encrypt:
    Enable: true
  Auth:
    Enable: true

#-------------------分害线
server:
  Listen: :4000
  HttpListen: :4009
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
	ClientCfg *ClientConfig `mapstructure:"client"`
	ServerCfg *ServerConfig `mapstructure:"server"`
	//OpenAuth  bool          `mapstructure:"openAuth"`
}

// ClientConfig read from a config file or cmd flags, or can be assgined from a server after got the server ack.
type ClientConfig struct {
	Registry   string `mapstructure:"server"`
	Listen     string `mapstructure:"listen"`
	HttpListen string `mapstructure:"httpListen"`
	ConsoleUrl string `mapstructure:"ConsoleUrl"`
	//Protocol   Protocol `mapstructure:"type"`
	Offset  int32   `mapstructure:"offset"`
	Encrypt Encrypt `mapstructure:"Encrypt"`
	Auth    Auth    `mapstructure:"Auth"`
	Relay   Relay   `mapstructure:"Relay"`
}

type Encrypt struct {
	Enable bool `mapstructure:"Enable"`
}

type Auth struct {
	Enable bool `mapstructure:"Enable"`
}

type Relay struct {
	Enable bool `mapstructure:"Enable"`
}

func (cfg *ClientConfig) EnableRelay() bool {
	return cfg.Relay.Enable
}

func (cfg *ClientConfig) EnableAuth() bool {
	return cfg.Auth.Enable
}

func (cfg *ClientConfig) EnableEncrypt() bool {
	return cfg.Encrypt.Enable
}

func (cfg *ClientConfig) HttpListenStr() string {
	return cfg.HttpListen
}

func (cfg *ClientConfig) HostUrl() string {
	return fmt.Sprintf("http://127.0.0.1%s", cfg.HttpListen)
}

func (cfg *ClientConfig) ControlUrl() string {
	return cfg.ConsoleUrl
}

func (cfg *ClientConfig) RegistryUrl() string {
	return cfg.Registry
}

func (cfg *ClientConfig) AuthEnable() bool {
	return cfg.Auth.Enable
}

type Redis struct {
	Enable   bool   `json:"enable"`
	Url      string `json:"host"`
	Port     int    `json:"port"`
	Username string `json:"username"`
	Password string `json:"password"`
}

type ServerConfig struct {
	Listen     string `mapstructure:"Listen"`
	HttpListen string `mapstructure:"HttpListen"`
	Redis      Redis  `json:redis`
	//Protocol   Protocol `mapstructure:"type"`
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
	viper.SetDefault("client.ConsoleUrl", "https://www.efvpn.com")
	viper.SetDefault("client.HttpListen", ":6662")
	viper.SetDefault("client.server", "www.efvpn.com:4000")
	viper.SetDefault("client.Relay.Enable", true)
	viper.SetDefault("client.Encrypt.Enable", true)
	viper.SetDefault("client.Auth.Enable", true)
	viper.SetDefault("client.Listen", ":3000")
	viper.SetDefault("server.Listen", ":4000")
	viper.SetDefault("server.HttpListen", ":4001")
	viper.SetDefault("server.redis.enable", false)
	if err = viper.ReadInConfig(); err != nil { // Handle errors reading the config file
		//if _, ok := err.(viper.ConfigFileNotFoundError); ok {
		//	viper.ReadConfig(bytes.NewBuffer(defaultYaml))
		//} else {
		//	return nil, errors.New("invalid config")
		//}
		return nil, err
	}

	if err = viper.UnmarshalExact(&config); err != nil {
		return nil, err
	}

	return

}
