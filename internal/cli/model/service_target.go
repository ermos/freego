package model

import (
	"net/url"
	"regexp"
)

type ServiceTarget struct {
	Url *url.URL
	Rgx *regexp.Regexp
}
