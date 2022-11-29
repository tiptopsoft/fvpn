package edge

import (
	"net"
)

func GetLocalMacAddr() string {
	// getMac
	ifaces, err := net.Interfaces()
	if err != nil {
		return ""
	}
	for _, v := range ifaces {
		if v.HardwareAddr == nil {
			continue
		}
		return v.HardwareAddr.String()
	}

	return ""
}
