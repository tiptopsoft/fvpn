package common

import (
	"errors"
	"github.com/spf13/pflag"
	"os/exec"
)

var (
	ErrNotImplemented = errors.New("not implement yet")
	ErrUnsupported    = errors.New("unsupported")
	ErrUnknow         = errors.New("unknown")
)

// ApplyFlags apply flag to struct
func ApplyFlags(fs *pflag.FlagSet) {

}

func ExecCommand(name string, commands ...string) *exec.Cmd {
	return exec.Command(name, commands...)
}
