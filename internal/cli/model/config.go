package model

import (
	"github.com/ermos/progo/internal/pkg/config"
)

type Config struct {
	ActiveDomains map[string]config.ActiveDomain `yaml:"activeDomains"`
}
