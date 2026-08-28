package validator

import (
	"testing"

	"go-backend-boilerplate/internal/shared/apperror"
	"go-backend-boilerplate/internal/shared/redact"
)

type sample struct {
	Email    string        `validate:"required,email"`
	Password redact.Secret `validate:"required,min=6"`
}

func TestStructValidatesSecret(t *testing.T) {
	err := Struct(sample{Email: "a@b.com", Password: "123"})
	if err == nil {
		t.Fatal("expected validation error for short password")
	}
	if e, ok := apperror.As(err); !ok || e.Status != 422 {
		t.Fatalf("expected 422 apperror, got %v", err)
	}

	if err := Struct(sample{Email: "a@b.com", Password: "123456"}); err != nil {
		t.Fatalf("expected no error for valid input, got %v", err)
	}

	if err := Struct(sample{Email: "bad", Password: "123456"}); err == nil {
		t.Fatal("expected validation error for bad email")
	}
}
