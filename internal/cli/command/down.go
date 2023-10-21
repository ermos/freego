package command

import (
	"github.com/ermos/progo/internal/cli/action"
	"github.com/ermos/progo/internal/pkg/config"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type Down struct {
	cfgFile string
}

func (d *Down) Flags(cmd *cobra.Command) {
	cmd.Flags().StringVarP(&d.cfgFile, "file", "f", "", "Progo configuration file (default is progo.yaml)")
}

func (d *Down) Execute(cmd *cobra.Command, args []string) error {
	c, err := action.GetAppConfig(d.cfgFile)
	if err != nil {
		return err
	}

	for id := range config.GetActiveDomainsFromLink(c.Link) {
		config.RemoveActiveDomain(id)
	}

	if err = viper.WriteConfig(); err != nil {
		return err
	}

	return nil
}
