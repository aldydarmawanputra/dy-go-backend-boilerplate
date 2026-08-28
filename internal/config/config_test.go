package config

import "testing"

func TestValidate(t *testing.T) {
	prodWeak := &Config{AppEnv: "production", JWTSecret: "change-me-in-production"}
	if err := prodWeak.Validate(); err == nil {
		t.Fatal("expected error for default JWT secret in production")
	}

	prodCORS := &Config{
		AppEnv: "production", JWTSecret: "a-sufficiently-long-secret",
		CORSAllowCredentials: true, CORSAllowOrigins: "*",
	}
	if err := prodCORS.Validate(); err == nil {
		t.Fatal("expected error for credentials + wildcard origin in production")
	}

	prodOK := &Config{AppEnv: "production", JWTSecret: "a-sufficiently-long-secret"}
	if err := prodOK.Validate(); err != nil {
		t.Fatalf("expected valid prod config, got %v", err)
	}

	dev := &Config{AppEnv: "development", JWTSecret: "dev"}
	if err := dev.Validate(); err != nil {
		t.Fatalf("expected valid dev config, got %v", err)
	}
}
