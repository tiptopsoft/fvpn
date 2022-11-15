package option

type RegConfig struct {
	Listen      string   `mapstructure:"listen"`
	Protocol    Protocol `mapstructure:"type"`
	OpenAuth    bool
	OpenEncrypt bool
}
