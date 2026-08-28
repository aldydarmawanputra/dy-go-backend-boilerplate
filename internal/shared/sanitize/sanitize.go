package sanitize

import (
	"strings"
	"unicode"
)

func String(s string) string {
	s = strings.TrimSpace(s)
	return strings.Map(func(r rune) rune {
		if unicode.IsControl(r) {
			return -1
		}
		return r
	}, s)
}

func Email(s string) string {
	return strings.ToLower(strings.TrimSpace(s))
}
