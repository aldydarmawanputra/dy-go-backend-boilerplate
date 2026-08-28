package redact

import (
	"encoding/json"
	"strings"
	"testing"
)

func TestSecretRedaction(t *testing.T) {
	s := Secret("supersecret")

	if got := s.Reveal(); got != "supersecret" {
		t.Fatalf("Reveal = %q", got)
	}
	if got := s.String(); got != mask {
		t.Fatalf("String = %q, want %q", got, mask)
	}

	payload := struct {
		Password Secret `json:"password"`
	}{Password: s}
	b, err := json.Marshal(payload)
	if err != nil {
		t.Fatal(err)
	}
	if strings.Contains(string(b), "supersecret") {
		t.Fatalf("marshaled output leaked the secret: %s", b)
	}
	if !strings.Contains(string(b), mask) {
		t.Fatalf("marshaled output not redacted: %s", b)
	}
}

func TestEmail(t *testing.T) {
	if got := Email("aldy@example.com"); got != "a***@example.com" {
		t.Fatalf("Email = %q", got)
	}
	if got := Email("notanemail"); got != mask {
		t.Fatalf("Email(no-at) = %q", got)
	}
}
