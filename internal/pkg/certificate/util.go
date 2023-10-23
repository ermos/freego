package certificate

import "strings"

func stringifyDomain(domain string) string {
	return strings.ReplaceAll(strings.ToLower(domain), "/", "-")
}
