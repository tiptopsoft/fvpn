package common

import (
	"errors"
	"fmt"
	"github.com/interstellar-cloud/star/option"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
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

func InitConfig() (config *option.Config, err error) {
	viper.SetConfigName("app") // name of config file (without extension)
	//viper.SetConfigType("yaml")           // REQUIRED if the config file does not have the extension in the name
	viper.AddConfigPath("/etc/star/")           // path to look for the config file in
	viper.AddConfigPath("$HOME/.star")          // call multiple times to add many search paths
	viper.AddConfigPath(".")                    // optionally look for config in the working directory
	if err = viper.ReadInConfig(); err != nil { // Handle errors reading the config file
		panic(fmt.Errorf("fatal error config file: %w", err))
	}

	if err = viper.UnmarshalExact(&config); err != nil {
		return nil, err
	}

	return

}
