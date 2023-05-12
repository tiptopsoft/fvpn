package cmd

import (
	"errors"
	"fmt"
	"github.com/spf13/cobra"
	"os"
	"path/filepath"
)

func logout() *cobra.Command {
	var opts loginOptions
	cmd := &cobra.Command{
		Use: "logout",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			return runLogout(&opts)
		},
	}

	return cmd
}

func runLogout(opts *loginOptions) error {
	//TOTO delete content in ~/.fvpn/config.json
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return err
	}

	path := filepath.Join(homeDir, "./fvpn/config.json")
	err = os.RemoveAll(path)
	if err != nil {
		return errors.New("logout failed")
	}

	fmt.Println("logout success")
	return nil
}
