package command

import (
	"fmt"
	"github.com/ermos/freego/internal/cli/model"
	"github.com/ermos/freego/internal/pkg/config"
	"github.com/ermos/freego/internal/pkg/util"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"regexp"
	"strings"
	"time"
)

type Serve struct {
	cfgFile string
}

func (d *Serve) Flags(cmd *cobra.Command) {
	cmd.Flags().StringVarP(&d.cfgFile, "file", "f", "", "freego configuration file (default is freego.yaml)")
}

func (d *Serve) Execute(cmd *cobra.Command, args []string) error {
	for {
		loadAt := time.Now().Unix()

		if err := config.Reload(); err != nil {
			panic(err)
		}

		var targets []model.ServiceTarget

		for _, domain := range config.GetActiveDomains() {
			rgx := regexp.QuoteMeta(domain.Domain)
			rgx = strings.ReplaceAll(rgx, `\*`, `.+`)
			rgx = fmt.Sprintf(`^%s.*$`, rgx)

			targets = append(targets, model.ServiceTarget{
				Url: d.urlParse(fmt.Sprintf("http://%s:%d", domain.Host, domain.Port)),
				Rgx: regexp.MustCompile(rgx),
			})
		}

		proxy := &httputil.ReverseProxy{
			Director: func(req *http.Request) {
				host := req.Host + req.URL.Path

				for _, target := range targets {
					if target.Rgx.MatchString(host) {
						log.Println(req.Host+req.URL.Path, " => ", target.Url.Host)
						req.URL.Scheme = target.Url.Scheme
						req.URL.Host = target.Url.Host
						return
					}
				}

			},
			ErrorHandler: func(writer http.ResponseWriter, request *http.Request, err error) {
				writer.Header().Set("Content-Type", "text/html; charset=utf-8")
				_, _ = writer.Write([]byte(d.errorPage(err)))
			},
		}

		server := &http.Server{
			Addr:    ":80",
			Handler: proxy,
		}

		serverErr := make(chan error, 1)

		go func() {
			err := server.ListenAndServe()
			serverErr <- err
		}()

		reload := make(chan bool, 1)

		go func() {
			for {
				time.Sleep(10 * time.Second)

				_ = config.Reload()

				if viper.GetInt64("lastupdate") > loadAt {
					reload <- true
					break
				}
			}
		}()

		select {
		case err := <-serverErr:
			panic(err)
		case <-reload:
			log.Println("Config updated ! freego is reloading...")
			if err := server.Close(); err != nil {
				panic(err)
			}
		}
	}

	return nil
}

func (Serve) urlParse(rawurl string) *url.URL {
	parsedURL, err := url.Parse(rawurl)
	if err != nil {
		panic(err)
	}
	return parsedURL
}

// @todo use template instead
func (Serve) errorPage(err error) string {
	var b []string

	b = append(b, "<!DOCTYPE html>")
	b = append(b, "<html>")
	b = append(b, "<head>")
	b = append(b, "<title>freego</title>")
	b = append(b, "<style>"+
		"html { "+
		"background-color: #162028; "+
		"color: #bbb;"+
		"font-size: 14px;"+
		"font-family: -apple-system,BlinkMacSystemFont,\"Segoe UI\",\"Noto Sans\",Helvetica,Arial,sans-serif,\"Apple Color Emoji\",\"Segoe UI Emoji\";"+
		"line-height: 1.5;"+
		"}"+
		"a { "+
		"color: #55fafa;"+
		"text-decoration: none;"+
		"font-weight: bold;"+
		"}"+
		"th, td { "+
		"padding: 5px 10px;"+
		"}"+
		"tr th:first-child, tr td:first-child { "+
		"padding-left: 0;"+
		"}"+
		"</style>")
	b = append(b, "</head>")
	b = append(b, "<body style=\"padding: 20px\">")

	// Freego ASCII art (Colossal Font)
	b = append(b, "<pre style=\"color: #55fafa;font-size: 10px;line-height: 1.1;\">")

	b = append(b, " .d888                                            <br>")
	b = append(b, "d88P\"                                             <br>")
	b = append(b, "888                                               <br>")
	b = append(b, "888888 888d888 .d88b.   .d88b.   .d88b.   .d88b.  <br>")
	b = append(b, "888    888P\"  d8P  Y8b d8P  Y8b d88P\"88b d88\"\"88b <br>")
	b = append(b, "888    888    88888888 88888888 888  888 888  888 <br>")
	b = append(b, "888    888    Y8b.     Y8b.     Y88b 888 Y88..88P <br>")
	b = append(b, "888    888     \"Y8888   \"Y8888   \"Y88888  \"Y88P\"  <br>")
	b = append(b, "                                     888           <br>")
	b = append(b, "                                Y8b d88P           <br>")
	b = append(b, "                                 \"Y88P\"            <br>")

	b = append(b, "</pre>")

	b = append(b, "<b>root</b> > echo $ERROR<br>")

	b = append(b, "<div style=\"color: red; font-weight: bold;\">")
	if err.Error() == "unsupported protocol scheme \"\"" {
		b = append(b, fmt.Sprintf("freego can't find port for the given url :("))
	} else {
		b = append(b, fmt.Sprintf("%s", err.Error()))
	}
	b = append(b, "</div>")

	b = append(b, "<b>root</b> > freego list<br>")

	b = append(b, "<table>")
	b = append(b, "<tr style=\"text-align: left\">")
	b = append(b, "<th>DOMAIN ID</th>")
	b = append(b, "<th>DOMAIN</th>")
	b = append(b, "<th>HOST</th>")
	b = append(b, "<th>PORT</th>")
	b = append(b, "<th>CREATED</th>")
	b = append(b, "</tr>")

	for domainId, domain := range config.GetActiveDomains() {
		b = append(b, "<tr>")
		b = append(b, fmt.Sprintf("<td>%s</td>", domainId))
		b = append(b, fmt.Sprintf("<td><a href=\"http://%s\">%s</a></td>", domain.Domain, domain.Domain))
		b = append(b, fmt.Sprintf("<td>%s</td>", domain.Host))
		b = append(b, fmt.Sprintf("<td>%d</td>", domain.Port))
		b = append(b, fmt.Sprintf("<td>%s</td>", util.FormatXTimeAgo(domain.CreatedAt, "never")))
		b = append(b, "</tr>")
	}

	b = append(b, "</body>")
	b = append(b, "</html>")

	return strings.Join(b, "")
}
