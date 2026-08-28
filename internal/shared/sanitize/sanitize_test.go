package sanitize

import "testing"

func TestString(t *testing.T) {
	cases := map[string]string{
		"  hello  ":        "hello",
		"a\x00b":           "ab",
		"line\x07break":    "linebreak",
		"  Aldy Darmawan ": "Aldy Darmawan",
	}
	for in, want := range cases {
		if got := String(in); got != want {
			t.Fatalf("String(%q) = %q, want %q", in, got, want)
		}
	}
}

func TestEmail(t *testing.T) {
	if got := Email("  Aldy@Example.COM "); got != "aldy@example.com" {
		t.Fatalf("Email = %q", got)
	}
}
