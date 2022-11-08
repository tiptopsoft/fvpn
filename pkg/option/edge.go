package option

// EdgeConfig read from a config file or cmd flags, or can be assgined from a registry after got the register ack.
type EdgeConfig struct {
	Registry string   `mapstructure:"registry"`
	Listen   string   `mapstructure:"listen"`
	TapName  string   `mapstructure:"tap"`
	TapIP    string   `mapstructure:"ip"`
	TapMask  string   `mapstructure:"mask"`
	MacAddr  string   `mapstructure:"mac"`
	Protocol Protocol `mapstructure:"type"`
}

func EdgeDefault() *EdgeConfig {
	return &EdgeConfig{}
}
