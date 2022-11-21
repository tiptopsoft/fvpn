package option

type RegConfig struct {
	Listen      string   `mapstructure:"listen"`
	HttpListen  string   `mapstructure:"httpListen"`
	Protocol    Protocol `mapstructure:"type"`
	OpenAuth    bool
	OpenEncrypt bool
}
