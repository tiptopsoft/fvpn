package common

import (
	"errors"
	"github.com/spf13/pflag"
	"os"
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

func ExecCommand(name string, commands ...string) error {
	cmd := exec.Command(name, commands...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}
