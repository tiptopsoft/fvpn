package option

// StarConfig conf for running a star up
type StarConfig struct {
	MoonIP string // service for moon server
	Port   int    // default port is 3000
	Server bool   // server or client, true: server
}
