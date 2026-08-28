package jwtutil

import (
	"testing"

	"github.com/golang-jwt/jwt/v5"
)

func newTestManager() *Manager {
	return New("test-secret", "test-issuer", "test-aud", 1)
}

func TestGenerateAndParse(t *testing.T) {
	m := newTestManager()

	token, err := m.Generate("user-123", []string{"admin", "user"})
	if err != nil {
		t.Fatalf("Generate: %v", err)
	}

	claims, err := m.Parse(token)
	if err != nil {
		t.Fatalf("Parse: %v", err)
	}
	if claims.Subject != "user-123" {
		t.Fatalf("subject = %q, want user-123", claims.Subject)
	}
	if len(claims.Roles) != 2 || claims.Roles[0] != "admin" {
		t.Fatalf("roles = %v, want [admin user]", claims.Roles)
	}
}

func TestGeneratedClaimsArePresent(t *testing.T) {
	m := newTestManager()
	token, _ := m.Generate("user-123", []string{"user"})

	claims := &Claims{}
	if _, _, err := jwt.NewParser().ParseUnverified(token, claims); err != nil {
		t.Fatalf("parse unverified: %v", err)
	}
	if claims.Issuer != "test-issuer" {
		t.Fatalf("iss = %q", claims.Issuer)
	}
	if len(claims.Audience) == 0 || claims.Audience[0] != "test-aud" {
		t.Fatalf("aud = %v", claims.Audience)
	}
	if claims.ID == "" {
		t.Fatal("jti (ID) is empty")
	}
	if claims.ExpiresAt == nil || claims.NotBefore == nil || claims.IssuedAt == nil {
		t.Fatal("exp/nbf/iat must be set")
	}
}

func TestParseRejectsWrongSecret(t *testing.T) {
	token, _ := New("secret-a", "iss", "aud", 1).Generate("user-123", nil)
	if _, err := New("secret-b", "iss", "aud", 1).Parse(token); err == nil {
		t.Fatal("expected error for wrong secret")
	}
}

func TestParseRejectsWrongIssuerOrAudience(t *testing.T) {
	token, _ := New("s", "iss-a", "aud-a", 1).Generate("user-123", nil)

	if _, err := New("s", "iss-b", "aud-a", 1).Parse(token); err == nil {
		t.Fatal("expected error for wrong issuer")
	}
	if _, err := New("s", "iss-a", "aud-b", 1).Parse(token); err == nil {
		t.Fatal("expected error for wrong audience")
	}
}
