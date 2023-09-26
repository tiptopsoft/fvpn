package device

import (
	"fmt"
	"github.com/olekukonko/tablewriter"
	"github.com/tiptopsoft/fvpn/pkg/util"
	"os"
)

func RunListNetworks(cfg *util.Config) error {
	logger.Debugf("start list networks")

	cm := NewManager(cfg.NodeCfg)
	resp, err := cm.ListNetworks()
	if err != nil {
		return err
	}

	var data [][]string
	if resp.List != nil {
		for index, id := range resp.List {
			data = append(data, []string{fmt.Sprintf("%d", index+1), id.NetworkId})
		}
	}

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Index", "Name"})

	table.AppendBulk(data)

	table.Render()
	return nil
}
