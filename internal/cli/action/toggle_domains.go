package action

import (
	"fmt"
	"github.com/ermos/freego/internal/cli/model"
	"github.com/ermos/freego/internal/pkg/term"
)

func ToggleDomains(c model.AppConfig, isDisabled bool) error {
	var domains []string

	for domain, content := range c.Domains {
		domains = append(domains, fmt.Sprintf("%s (%s:%d)", domain, content.Host, content.Port))
	}

	return term.ToggleDomain(domains, isDisabled)
}
