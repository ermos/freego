package main

import (
	"fmt"
	"github.com/ermos/progo/internal/pkg/config"
	"github.com/spf13/viper"
	"net/http"
	"net/http/httputil"
	"net/url"
	"regexp"
	"strings"
	"time"
)

type Target struct {
	url *url.URL
	rgx *regexp.Regexp
}

func main() {
	for {
		loadAt := time.Now().Unix()

		if err := config.Init(); err != nil {
			panic(err)
		}

		var targets []Target

		for _, domain := range config.GetActiveDomains() {
			rgx := regexp.QuoteMeta(domain.Domain)
			rgx = strings.ReplaceAll(rgx, `\*`, `.+`)
			rgx = fmt.Sprintf(`^%s.*$`, rgx)

			targets = append(targets, Target{
				url: urlParse(fmt.Sprintf("http://%s:%d", domain.Host, domain.Port)),
				rgx: regexp.MustCompile(rgx),
			})
		}

		proxy := &httputil.ReverseProxy{
			Director: func(req *http.Request) {
				host := req.Host + req.URL.Path

				for _, target := range targets {
					if target.rgx.MatchString(host) {
						req.URL.Scheme = target.url.Scheme
						req.URL.Host = target.url.Host
						return
					}
				}

			},
			ErrorHandler: func(writer http.ResponseWriter, request *http.Request, err error) {
				var message string

				// Progo ASCII art
				message += "    ____                       \n   / __ \\_________  ____ _____ \n  / /_/ / ___/ __ \\/ __ `/ __ \\\n / ____/ /  / /_/ / /_/ / /_/ /\n/_/   /_/   \\____/\\__, /\\____/ \n                 /____/        \n\n"

				if err.Error() == "unsupported protocol scheme \"\"" {
					message += fmt.Sprintf("Progo can't find port for the given url :(\n\n")
				} else {
					message += fmt.Sprintf("%s\n\n", err.Error())
				}

				message += "Active domains:\n"
				for _, domain := range config.GetActiveDomains() {
					message += fmt.Sprintf("  - %s => %s:%d\n", domain.Domain, domain.Host, domain.Port)
				}

				http.Error(writer, message, http.StatusNotFound)
			},
		}

		server := &http.Server{
			Addr:    ":80",
			Handler: proxy,
		}

		go server.ListenAndServe()

		for {
			time.Sleep(10 * time.Second)

			_ = config.Reload()

			if viper.GetInt64("lastupdate") > loadAt {
				break
			}
		}

		if err := server.Close(); err != nil {
			panic(err)
		}
	}
}

func urlParse(rawurl string) *url.URL {
	parsedURL, err := url.Parse(rawurl)
	if err != nil {
		panic(err)
	}
	return parsedURL
}
