package redact

import (
	"log/slog"
	"strings"
)

const mask = "[REDACTED]"

type Secret string

func (s Secret) Reveal() string { return string(s) }

func (Secret) String() string { return mask }

func (Secret) MarshalJSON() ([]byte, error) { return []byte(`"` + mask + `"`), nil }

func (Secret) LogValue() slog.Value { return slog.StringValue(mask) }

func String(s string) string {
	if s == "" {
		return ""
	}
	return mask
}

func Email(s string) string {
	at := strings.IndexByte(s, '@')
	if at <= 0 {
		return mask
	}
	local := s[:at]
	if len(local) <= 1 {
		return "*" + s[at:]
	}
	return string(local[0]) + "***" + s[at:]
}

func Tail(s string, keep int) string {
	if keep <= 0 || len(s) <= keep {
		return mask
	}
	return mask + s[len(s)-keep:]
}
