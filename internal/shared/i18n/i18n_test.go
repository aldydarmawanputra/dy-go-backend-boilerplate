package i18n

import "testing"

func TestResolve(t *testing.T) {
	cases := map[string]string{
		"id-ID,id;q=0.9,en;q=0.8": "id",
		"en-US,en;q=0.9":          "en",
		"fr-FR,fr;q=0.9":          "en", // unsupported -> default
		"":                        "en",
		"ID":                      "id",
	}
	for in, want := range cases {
		if got := Resolve(in); got != want {
			t.Fatalf("Resolve(%q) = %q, want %q", in, got, want)
		}
	}
}

func TestT(t *testing.T) {
	if got := T("id", "CONFLICT", "fallback"); got != "data sudah ada" {
		t.Fatalf("T id CONFLICT = %q", got)
	}
	if got := T("en", "NOT_FOUND", "fallback"); got != "resource not found" {
		t.Fatalf("T en NOT_FOUND = %q", got)
	}
	if got := T("id", "UNKNOWN_CODE", "fallback"); got != "fallback" {
		t.Fatalf("T unknown = %q, want fallback", got)
	}
}
