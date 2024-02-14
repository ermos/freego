package config

import (
	"fmt"
	"github.com/spf13/viper"
	"net/url"
	"time"
)

type ActiveDomain struct {
	Link      string    `yaml:"link"`
	Domain    string    `yaml:"domain"`
	Host      string    `yaml:"host"`
	Port      int       `yaml:"port"`
	CreatedAt time.Time `yaml:"createdat"`
}

func (ad ActiveDomain) GetBaseDomain() string {
	u, err := url.Parse("https://" + ad.Domain)
	if err != nil {
		return ""
	}

	return u.Host
}

func GetActiveDomains() map[string]ActiveDomain {
	var activeDomains map[string]ActiveDomain

	err := viper.UnmarshalKey("active_domains", &activeDomains)
	if err != nil {
		panic(err)
	}

	return activeDomains
}

func GetActiveDomainsFromLink(link string) map[string]ActiveDomain {
	activeDomains := make(map[string]ActiveDomain)

	for id, ad := range GetActiveDomains() {
		if ad.Link == link {
			activeDomains[id] = ad
		}
	}

	return activeDomains
}

func AddActiveDomain(id string, ad ActiveDomain) {
	viper.Set(fmt.Sprintf("active_domains.%s", id), ad)
}

func RemoveActiveDomain(id string) {
	activeDomains := make(map[string]ActiveDomain)

	// viper doesn't support delete key from map
	for adId, ad := range GetActiveDomains() {
		if adId != id {
			activeDomains[adId] = ad
		}
	}

	viper.Set("active_domains", activeDomains)
}
