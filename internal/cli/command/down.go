package command

import (
	"github.com/ermos/freego/internal/cli/action"
	"github.com/spf13/cobra"
)

type Down struct {
	cfgFile string
}

func (d *Down) Flags(cmd *cobra.Command) {
	cmd.Flags().StringVarP(&d.cfgFile, "file", "f", "", "freego configuration file (default is freego.yaml)")
}

func (d *Down) Execute(cmd *cobra.Command, args []string) error {
	var toggleDomains []action.ToggleDomain

	c, err := action.GetAppConfig(d.cfgFile)
	if err != nil {
		return err
	}

	err = action.RemoveHostsFromLink(c.Link)
	if err != nil {
		return err
	}

	return action.ToggleDomains(toggleDomains)
}
