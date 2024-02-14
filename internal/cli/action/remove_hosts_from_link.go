package action

import (
	"github.com/ermos/freego/internal/cli/model"
	"github.com/ermos/freego/internal/pkg/config"
	"github.com/spf13/viper"
)

func RemoveHostsFromLink(link string) error {
	var toggleDomains []ToggleDomain

	for id, domain := range config.GetActiveDomainsFromLink(link) {
		config.RemoveActiveDomain(id)

		toggleDomains = append(toggleDomains, ToggleDomain{
			DomainName: domain.Domain,
			Domain: model.Domain{
				Host: domain.Host,
				Port: domain.Port,
			},
			Status: "Removed",
		})
	}

	SetLastUpdate()

	return viper.WriteConfig()
}
