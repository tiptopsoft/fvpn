package option

import "sync"

// StarConfig read from a config file or cmd flags, or can be assgined from a registry after got the registry ack.
type StarConfig struct {
	Registry string   `mapstructure:"registry"`
	Listen   string   `mapstructure:"listen"`
	TapName  string   `mapstructure:"tap"`
	TapIP    string   `mapstructure:"ip"`
	TapMask  string   `mapstructure:"mask"`
	MacAddr  string   `mapstructure:"mac"`
	Protocol Protocol `mapstructure:"type"`
}

func Default() *StarConfig {
	return &StarConfig{}
}

var (
	AddrMap sync.Map
)
