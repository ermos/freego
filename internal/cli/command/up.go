package command

import (
	"github.com/ermos/freego/internal/cli/action"
	"github.com/ermos/freego/internal/pkg/certificate"
	"github.com/ermos/freego/internal/pkg/config"
	"github.com/google/uuid"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"strings"
	"time"
)

type Up struct {
	cfgFile string
}

func (u *Up) Flags(cmd *cobra.Command) {
	cmd.Flags().StringVarP(&u.cfgFile, "file", "f", "", "freego configuration file (default is freego.yaml)")
}

func (u *Up) Execute(cmd *cobra.Command, args []string) error {
	var toggleDomains []action.ToggleDomain

	c, err := action.GetAppConfig(u.cfgFile)
	if err != nil {
		return err
	}

	linkedDomains := config.GetActiveDomainsFromLink(c.Link)

	for domain, content := range c.Domains {
		var id string
		var toggleDomain action.ToggleDomain

		for gdId, gd := range linkedDomains {
			if gd.Domain == domain {
				id = gdId
				break
			}
		}

		if id == "" {
			var idUuid uuid.UUID

			idUuid, err = uuid.NewRandom()
			if err != nil {
				return err
			}

			id = strings.ReplaceAll(idUuid.String(), "-", "")

			toggleDomain.Status = "Added"
		} else {
			toggleDomain.Status = "Updated"
		}

		config.AddActiveDomain(id, config.ActiveDomain{
			Domain:    domain,
			Host:      content.Host,
			Port:      content.Port,
			Link:      c.Link,
			CreatedAt: time.Now(),
		})

		toggleDomain.DomainName = domain
		toggleDomain.Domain = content

		toggleDomains = append(toggleDomains, toggleDomain)
	}

	for gdId, gd := range linkedDomains {
		isDelete := true

		for domain := range c.Domains {
			if gd.Domain == domain {
				isDelete = false
				break
			}
		}

		if isDelete {
			config.RemoveActiveDomain(gdId)
		}
	}

	action.SetLastUpdate()

	for domain := range c.Domains {
		if err = certificate.Generate(domain); err != nil {
			return err
		}
	}

	if err = viper.WriteConfig(); err != nil {
		return err
	}

	return action.ToggleDomains(toggleDomains)
}
