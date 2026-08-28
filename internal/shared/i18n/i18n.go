package i18n

import "strings"

const Default = "en"

var catalog = map[string]map[string]string{
	"en": {
		"BAD_REQUEST":       "invalid request",
		"UNAUTHORIZED":      "unauthorized",
		"FORBIDDEN":         "you are not allowed to do this",
		"NOT_FOUND":         "resource not found",
		"CONFLICT":          "resource already exists",
		"VALIDATION_ERROR":  "validation failed",
		"TOO_MANY_REQUESTS": "too many requests",
		"INTERNAL":          "internal server error",
	},
	"id": {
		"BAD_REQUEST":       "permintaan tidak valid",
		"UNAUTHORIZED":      "tidak terautentikasi",
		"FORBIDDEN":         "kamu tidak diizinkan melakukan ini",
		"NOT_FOUND":         "data tidak ditemukan",
		"CONFLICT":          "data sudah ada",
		"VALIDATION_ERROR":  "validasi gagal",
		"TOO_MANY_REQUESTS": "terlalu banyak permintaan",
		"INTERNAL":          "terjadi kesalahan pada server",
	},
}

func Supported(locale string) bool {
	_, ok := catalog[locale]
	return ok
}

// T returns the translation for key in locale, or fallback when unavailable.
func T(locale, key, fallback string) string {
	if m, ok := catalog[locale]; ok {
		if v, ok := m[key]; ok {
			return v
		}
	}
	return fallback
}

// Resolve picks the best supported locale from an Accept-Language header value
// like "id-ID,id;q=0.9,en;q=0.8", falling back to Default.
func Resolve(acceptLanguage string) string {
	for _, part := range strings.Split(acceptLanguage, ",") {
		tag := strings.TrimSpace(part)
		if i := strings.IndexByte(tag, ';'); i >= 0 {
			tag = tag[:i]
		}
		tag = strings.ToLower(strings.TrimSpace(tag))
		if i := strings.IndexByte(tag, '-'); i >= 0 {
			tag = tag[:i]
		}
		if Supported(tag) {
			return tag
		}
	}
	return Default
}
