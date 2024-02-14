package action

import (
	"fmt"
	"github.com/ermos/freego/internal/cli/model"
	"github.com/ermos/freego/internal/pkg/term"
)

type ToggleDomain struct {
	DomainName string
	Domain     model.Domain
	Status     string
}

func (td ToggleDomain) ToString() string {
	return fmt.Sprintf("%s (%s:%d) [%s]", td.DomainName, td.Domain.Host, td.Domain.Port, td.Status)
}

func ToggleDomains(domains []ToggleDomain) error {
	var list []string

	for _, domain := range domains {
		list = append(list, domain.ToString())
	}

	return term.ToggleDomain(list)
}
