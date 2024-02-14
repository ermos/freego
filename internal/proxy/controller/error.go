package controller

import (
	"fmt"
	"github.com/ermos/freego/internal/pkg/config"
	"github.com/ermos/freego/internal/pkg/util"
	"net/http"
	"strings"
)

type ErrorController struct{}

func (ec ErrorController) Handler() func(http.ResponseWriter, *http.Request, error) {
	return func(w http.ResponseWriter, r *http.Request, err error) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")

		_, writerErr := w.Write([]byte(ec.errorPage(err)))
		if writerErr != nil {
			panic(writerErr)
		}
	}
}

// @todo use template instead
func (ec ErrorController) errorPage(err error) string {
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
