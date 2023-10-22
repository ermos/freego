package command

import (
	"github.com/ermos/freego/internal/cli/action"
	"github.com/ermos/freego/internal/pkg/config"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type Down struct {
	cfgFile string
}

func (d *Down) Flags(cmd *cobra.Command) {
	cmd.Flags().StringVarP(&d.cfgFile, "file", "f", "", "freego configuration file (default is freego.yaml)")
}

func (d *Down) Execute(cmd *cobra.Command, args []string) error {
	c, err := action.GetAppConfig(d.cfgFile)
	if err != nil {
		return err
	}

	for id := range config.GetActiveDomainsFromLink(c.Link) {
		config.RemoveActiveDomain(id)
	}

	action.SetLastUpdate()

	if err = viper.WriteConfig(); err != nil {
		return err
	}

	return nil
}
