package option

const (
	TCP = iota
	UDP
)

// StarConfig conf for running a star up
type StarConfig struct {
	Star
	MoonIP string // service for moon server
	Port   int    // default port is 3000
	Server bool   // server or client, true: server
	Mac    string // like "07:00:10:24:55:42"
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
}

type Mysql struct {
	User     string `mapstructure:"user"`
	Url      string `mapstructure:"url"`
	Password string `mapstructure:"password"`
	Name     string `mapstructure:"name"`
}
