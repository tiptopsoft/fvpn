package util

import "github.com/sevlyar/go-daemon"

var (
	daemonLogFile = "fvpn.log"
	daemonPidFile = "fvpn.pid"
)

func GetDaemon() *daemon.Context {
	return &daemon.Context{
		PidFileName: daemonPidFile,
		PidFilePerm: 0644,
		LogFileName: daemonLogFile,
		LogFilePerm: 0640,
		WorkDir:     "./",
		Umask:       027,
	}

}
