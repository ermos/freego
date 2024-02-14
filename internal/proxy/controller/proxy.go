package controller

import (
	"fmt"
	"github.com/ermos/freego/internal/cli/model"
	"github.com/ermos/freego/internal/pkg/config"
	"log"
	"net/http"
	"net/url"
	"regexp"
	"strings"
)

type ProxyController struct{}

func (pc ProxyController) Handler() func(req *http.Request) {
	var targets []model.ServiceTarget

	for _, domain := range config.GetActiveDomains() {
		rgx := regexp.QuoteMeta(domain.Domain)
		rgx = strings.ReplaceAll(rgx, `\*`, `[^\/]+`)
		rgx = fmt.Sprintf(`^%s(.*)$`, rgx)

		targets = append(targets, model.ServiceTarget{
			Url: pc.urlParse(fmt.Sprintf("https://%s:%d", domain.Host, domain.Port)),
			Rgx: regexp.MustCompile(rgx),
		})
	}

	return func(req *http.Request) {
		host := req.Host + req.URL.Path

		for _, target := range targets {
			if target.Rgx.MatchString(host) {
				log.Println(req.Host+req.URL.Path, " => ", target.Url.Host)
				req.URL.Scheme = target.Url.Scheme
				req.URL.Host = target.Url.Host
				req.URL.Path = target.Rgx.FindStringSubmatch(host)[1]
				return
			}
		}

	}
}

func (ProxyController) urlParse(rawurl string) *url.URL {
	parsedURL, err := url.Parse(rawurl)
	if err != nil {
		panic(err)
	}
	return parsedURL
}
