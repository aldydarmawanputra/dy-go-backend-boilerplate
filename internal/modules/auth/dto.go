package auth

import "go-backend-boilerplate/internal/shared/redact"

type RegisterRequest struct {
	Email    string        `json:"email" validate:"required,email"`
	Name     string        `json:"name" validate:"required,min=2,max=100"`
	Password redact.Secret `json:"password" validate:"required,min=6,max=72"`
}

type LoginRequest struct {
	Email    string        `json:"email" validate:"required,email"`
	Password redact.Secret `json:"password" validate:"required"`
}

type TokenResponse struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
}
