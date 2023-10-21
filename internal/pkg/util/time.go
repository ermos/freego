package util

import (
	"fmt"
	"time"
)

func FormatXTimeAgo(t time.Time, neverText string) string {
	if t == (time.Time{}) {
		return neverText
	}

	duration := time.Since(t)

	years := int(duration.Hours() / 24 / 365)
	months := int(duration.Hours()/24/30) % 12
	days := int(duration.Hours()/24) % 30
	hours := int(duration.Hours()) % 24
	minutes := int(duration.Minutes()) % 60

	var result string

	switch {
	case years >= 1:
		result = fmt.Sprintf("%d year%s ago", years, pluralize(years))
		break
	case months >= 1:
		result = fmt.Sprintf("%d month%s ago", months, pluralize(months))
		break
	case days >= 1:
		result = fmt.Sprintf("%d day%s ago", days, pluralize(days))
		break
	case hours >= 1:
		result = fmt.Sprintf("%d hour%s ago", hours, pluralize(hours))
		break
	case minutes >= 1:
		result = fmt.Sprintf("%d minute%s ago", minutes, pluralize(minutes))
		break
	default:
		result = "less than a minute ago"
	}

	return result
}

func pluralize(count int) string {
	if count == 1 {
		return ""
	} else {
		return "s"
	}
}
