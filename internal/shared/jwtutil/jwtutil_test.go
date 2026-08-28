package jwtutil

import "testing"

func TestGenerateAndParse(t *testing.T) {
	m := New("test-secret", 1)

	token, err := m.Generate("user-123")
	if err != nil {
		t.Fatalf("Generate returned error: %v", err)
	}

	userID, err := m.Parse(token)
	if err != nil {
		t.Fatalf("Parse returned error: %v", err)
	}
	if userID != "user-123" {
		t.Fatalf("expected subject user-123, got %q", userID)
	}
}

func TestParseRejectsWrongSecret(t *testing.T) {
	token, err := New("secret-a", 1).Generate("user-123")
	if err != nil {
		t.Fatalf("Generate returned error: %v", err)
	}

	if _, err := New("secret-b", 1).Parse(token); err == nil {
		t.Fatal("expected error when parsing with a different secret, got nil")
	}
}
